// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"

	"github.com/syhpoon/learnfs/pkg/device"
	"github.com/syhpoon/learnfs/pkg/lib"

	"github.com/stretchr/testify/require"
)

type fuzzFile struct {
	name string
	size int
	md5  string
}

type fuzzState struct {
	dirLevels    uint
	files        []*fuzzFile
	dirs         []string
	maxNumInodes int
	fsObj        *Filesystem
}

func (st *fuzzState) isFull() bool {
	return len(st.dirs)+len(st.files) >= st.maxNumInodes
}

func FuzzFs(f *testing.F) {
	// seed, number of dir levels
	f.Add(int64(0), uint(10))
	f.Fuzz(fuzzTest)
}

func fuzzTest(t *testing.T, seed int64, dirLevels uint) {
	ctx := context.Background()
	rnd := rand.New(rand.NewSource(seed))

	// Select a device size from 100 to 300 MiB
	devSize := (rnd.Intn(200) + 100) * int(math.Pow(1024, 2))
	dev := device.NewMemory(ctx, uint(devSize))

	err := Create(dev)
	require.NoError(t, err)

	fsObj, err := Load(dev)
	require.NoError(t, err)

	state := &fuzzState{
		dirLevels:    dirLevels % 5,
		maxNumInodes: int(fsObj.superblock.NumInodes),
		fsObj:        fsObj,
	}

	generateLevel(t, "/", 1, rnd, state)

	sort.Slice(state.files, func(i, j int) bool {
		return state.files[i].name < state.files[j].name
	})
	sort.Strings(state.dirs)

	inodeMap := map[string]InodePtr{"/": 1}
	fileData := map[InodePtr][]byte{}

	uid := uint32(1000)
	gid := uint32(1500)

	// Create dirs
	for _, dir := range state.dirs {
		parentInode := InodePtr(1) // root dir inode
		var parts []string

		for _, part := range strings.Split(dir[1:], "/") {
			parts = append(parts, part)
			var ptr InodePtr

			partsStr := filepath.Join(parts...)

			if ino, ok := inodeMap[partsStr]; ok {
				ptr = ino
			} else {
				ino, err := fsObj.CreateDirectory(parentInode, part,
					0x755|syscall.S_IFDIR, 0, uid, gid)
				require.NoError(t, err, "should create directory '%s' (%s)",
					part, strings.Join(parts, "/"))
				ptr = ino.ptr
			}

			inodeMap[partsStr] = ptr
			parentInode = ptr
		}
	}

	dirFiles := map[string][]*fuzzFile{}

	// Create files
	for i := range state.files {
		file := state.files[i]
		dir, f := filepath.Split(file.name)
		if len(dir) > 1 {
			dir = strings.Trim(dir, "/")
		}

		dirIno, ok := inodeMap[dir]
		require.Truef(t, ok, "containing directory %s must exist", dir)

		ino, err := fsObj.CreateFile(dirIno, f, 0x644, 0, uid, gid)
		require.NoError(t, err)

		data := make([]byte, file.size)
		rnd.Read(data)

		_, err = fsObj.Write(ino.ptr, 0, data)
		require.NoError(t, err)

		hash := md5.Sum(data)
		file.md5 = hex.EncodeToString(hash[:])

		fileData[ino.ptr] = data
		dirFiles[dir] = append(dirFiles[dir], file)
	}

	fsObj.flusher.flushAll()
	fsObj.flusher.flushBitmaps()

	verifyHier := func(fsObj *Filesystem) {
		for _, dir := range state.dirs {
			dir := strings.Trim(dir, "/")
			dirInode, ok := inodeMap[dir]
			require.Truef(t, ok, "dir %s inode must be present in cache: %v", dir, inodeMap)

			dirObj, err := fsObj.GetDir(dirInode)
			require.NoError(t, err, "must get dir %s (%d)", dir, dirInode)

			for _, file := range dirFiles[dir] {
				f := filepath.Base(file.name)
				fileInode, err := dirObj.GetEntry(f)
				require.NoError(t, err, "must get file %s from dir %s", f, dir)
				require.NotZero(t, fileInode)
			}
		}
	}

	// Verify hierarchy
	verifyHier(fsObj)

	// Load fs anew
	fsObj2, err := Load(dev)
	require.NoError(t, err)

	verifyHier(fsObj2)

	rnd.Shuffle(len(state.files), func(i, j int) {
		state.files[i], state.files[j] = state.files[j], state.files[i]
	})

	// Stat and read all the files
	for _, file := range state.files {
		dir, f := filepath.Split(file.name)
		if len(dir) > 1 {
			dir = strings.Trim(dir, "/")
		}

		dirIno, ok := inodeMap[dir]
		require.Truef(t, ok, "containing directory %s must exist", dir)

		fileIno, err := fsObj.Lookup(dirIno, f)
		require.NoError(t, err, "must lookup file %s from dir %s", f, dir)

		require.Equal(t, fileIno.Uid, uid)
		require.Equal(t, fileIno.Gid, gid)
		require.Equal(t, fileIno.Size, uint32(file.size))

		data, err := fsObj.Read(fileIno.ptr, 0, int(fileIno.Size))
		require.NoError(t, err, "must read file %s in dir %s", f, dir)

		hash := md5.Sum(data)
		md5Sum := hex.EncodeToString(hash[:])
		require.Equal(t, file.md5, md5Sum, "must equal file %s content md5", f)
	}

	// Delete files
	for _, file := range state.files {
		dir, f := filepath.Split(file.name)
		if len(dir) > 1 {
			dir = strings.Trim(dir, "/")
		}

		dirIno, ok := inodeMap[dir]
		require.Truef(t, ok, "containing directory %s must exist", dir)

		err := fsObj.RemoveFile(dirIno, f)
		require.NoError(t, err, "must delete file %s from dir %s", f, dir)

		ino, err := fsObj.Lookup(dirIno, f)
		require.Nil(t, ino)
		require.True(t, errors.Is(err, ErrorNotFound))
	}

	// Delete dirs
	sort.Slice(state.dirs, func(i, j int) bool {
		return state.dirs[i] > state.dirs[j]
	})

	for _, dir := range state.dirs {
		parent, child := filepath.Split(dir)

		if parent == "/" {
			continue
		}

		parent = strings.Trim(parent, "/")

		dirIno, ok := inodeMap[parent]
		require.Truef(t, ok, "containing directory %s must exist", parent)

		err := fsObj.RemoveDirectory(dirIno, child)
		require.NoError(t, err, "must delete child dir %s from dir %s", child, parent)

		ino, err := fsObj.Lookup(dirIno, child)
		require.Nil(t, ino)
		require.True(t, errors.Is(err, ErrorNotFound))
	}
}

func generateLevel(t *testing.T, dir string, level uint, rnd *rand.Rand, state *fuzzState) {
	randName := func(i int) string {
		s := lib.RandomBytes(10, rnd) + fmt.Sprintf(".%d", i)

		return s
	}

	if level > state.dirLevels {
		return
	}

	// Number of files on this dir level
	numFiles := rnd.Intn(20)

	for f := 0; f < numFiles; f++ {
		// 1 byte to 500K
		fileSize := rnd.Intn(512000) + 1

		fileBlocks := uint32(fileSize) / (state.fsObj.blockSize)
		freeBlocks := state.fsObj.blockAllocator.GetBitmap().numFreeBits()

		if freeBlocks < fileBlocks {
			return
		}

		state.files = append(state.files,
			&fuzzFile{
				name: filepath.Join(dir, randName(f)),
				size: fileSize,
			})

		if state.isFull() {
			return
		}
	}

	// Number of dirs on this dir level
	numDirs := rnd.Intn(4)

	for d := 0; d < numDirs; d++ {
		newDir := filepath.Join(dir, randName(d))
		state.dirs = append(state.dirs, newDir)

		if state.isFull() {
			return
		}

		generateLevel(t, newDir, level+1, rnd, state)
	}
}

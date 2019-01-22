package store

import (
	"testing"

	"flamingo.me/flamingo/v3/core/auth/domain"
	"github.com/stretchr/testify/suite"
	"github.com/zemirco/memorystore"
)

type (
	MemoryTestSuite struct {
		suite.Suite

		memory *Memory

		store *memorystore.MemoryStore
	}
)

func TestMemoryTestSuite(t *testing.T) {
	suite.Run(t, &MemoryTestSuite{})
}

func (t *MemoryTestSuite) SetupTest() {
	t.store = memorystore.NewMemoryStore()
	t.memory = &Memory{}
	t.memory.Inject(t.store)
}

func (t *MemoryTestSuite) TearDownTest() {
	t.store = nil
	t.memory = nil
}

func (t *MemoryTestSuite) TestDestroySessionsForUser_Error() {
	user := domain.User{
		Sub: "sub",
	}

	all := t.store.GetAll()
	key := t.memory.getAllHashesKey(user)
	all[key] = "wrong"

	t.Error(t.memory.DestroySessionsForUser(user))
}

func (t *MemoryTestSuite) TestDestroySessionsForUser_SuccessEmptyMap() {
	user := domain.User{
		Sub: "sub",
	}

	t.NoError(t.memory.DestroySessionsForUser(user))

	all := t.store.GetAll()
	key := t.memory.getAllHashesKey(user)
	_, ok := all[key]
	t.False(ok)
}

func (t *MemoryTestSuite) TestDestroySessionsForUser_SuccessListInMap() {
	user := domain.User{
		Sub: "sub",
	}

	all := t.store.GetAll()
	key := t.memory.getAllHashesKey(user)
	all["session1"] = "session1"
	all["session2"] = "session2"
	all[key] = `["session1", "session2"]`

	_, ok := all["session1"]
	t.True(ok)
	_, ok = all["session2"]
	t.True(ok)
	_, ok = all[key]
	t.True(ok)

	t.NoError(t.memory.DestroySessionsForUser(user))

	_, ok = all["session1"]
	t.False(ok)
	_, ok = all["session2"]
	t.False(ok)
	_, ok = all[key]
	t.False(ok)
}

func (t *MemoryTestSuite) TestGetHashByUser_Error() {
	user := domain.User{
		Sub: "sub",
	}

	hash, err := t.memory.GetHashByUser(user)
	t.Error(err)
	t.Equal("", hash)
}

func (t *MemoryTestSuite) TestGetHashByUser_Success() {
	user := domain.User{
		Sub: "sub",
	}

	all := t.store.GetAll()
	key := t.memory.getHashKey(user)
	all[key] = "hash"

	hash, err := t.memory.GetHashByUser(user)
	t.NoError(err)
	t.Equal("hash", hash)
}

func (t *MemoryTestSuite) TestSetHashAndSessionIdForUser_Error() {
	user := domain.User{
		Sub: "sub",
	}

	all := t.store.GetAll()
	key := t.memory.getAllHashesKey(user)
	all[key] = "wrong"

	t.Error(t.memory.SetHashAndSessionIdForUser(user, "hash", "id"))
}

func (t *MemoryTestSuite) TestSetHashAndSessionIdForUser_EmptyMap() {
	user := domain.User{
		Sub: "sub",
	}

	all := t.store.GetAll()
	keyAll := t.memory.getAllHashesKey(user)
	key := t.memory.getHashKey(user)

	_, ok := all[keyAll]
	t.False(ok)
	_, ok = all[key]
	t.False(ok)

	t.NoError(t.memory.SetHashAndSessionIdForUser(user, "hash1", "id1"))

	ids, ok := all[keyAll]
	t.True(ok)
	t.Equal(`["id1"]`, ids)

	hash, ok := all[key]
	t.True(ok)
	t.Equal("hash1", hash)
}

func (t *MemoryTestSuite) TestSetHashAndSessionIdForUser_Replace() {
	user := domain.User{
		Sub: "sub",
	}

	all := t.store.GetAll()
	keyAll := t.memory.getAllHashesKey(user)
	key := t.memory.getHashKey(user)
	all[keyAll] = `["id2"]`
	all[key] = "id2"

	_, ok := all[keyAll]
	t.True(ok)
	_, ok = all[key]
	t.True(ok)

	t.NoError(t.memory.SetHashAndSessionIdForUser(user, "hash1", "id1"))

	ids, ok := all[keyAll]
	t.True(ok)
	t.Equal(`["id2","id1"]`, ids)

	hash, ok := all[key]
	t.True(ok)
	t.Equal("hash1", hash)
}

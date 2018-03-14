package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockAddParent(t *testing.T) {
	var p, c Block
	assert.False(t, p.IsParentOf(c))
	assert.False(t, c.IsParentOf(p))

	// c.Height is 0
	err := c.AddParent(p)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "height")
	}
	c.Height = 100
	err = c.AddParent(p)
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "height")
	}

	c.Height = p.Height + 1
	assert.NoError(t, c.AddParent(p))
	assert.True(t, p.IsParentOf(c))
	assert.False(t, c.IsParentOf(p))
}

func TestDecodeBlock(t *testing.T) {
	assert := assert.New(t)

	m1 := NewMessage(Address("from"), Address("to"), big.NewInt(10), "hello", []byte("cat"))
	m2 := NewMessage(Address("from"), Address("other"), big.NewInt(2), "yes", []byte("dog"))

	m1Cid, err := m1.Cid()
	assert.NoError(err)
	m2Cid, err := m2.Cid()
	assert.NoError(err)

	c1, err := cidFromString("a")
	assert.NoError(err)
	c2, err := cidFromString("b")
	assert.NoError(err)

	before := &Block{
		Parent:    c1,
		Height:    2,
		Messages:  []*Message{m1, m2},
		StateRoot: c2,
		MessageReceipts: []*MessageReceipt{
			NewMessageReceipt(m1Cid, 1, "", []byte{1, 2}),
			NewMessageReceipt(m2Cid, 1, "", []byte{1, 2}),
		},
	}

	after, err := DecodeBlock(before.ToNode().RawData())
	assert.NoError(err)
	assert.Equal(after.Cid(), before.Cid())
	assert.Equal(after, before)
}
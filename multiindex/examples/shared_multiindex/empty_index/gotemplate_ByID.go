// Code generated by gotemplate. DO NOT EDIT.

package empty_index

import (
	"fmt"
	"unsafe"

	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/multiindex"
	. "github.com/eosspark/eos-go/common/container/offsetptr"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply,Allocator)

// OrderedIndex holds elements of the red-black tree
type ByID struct {
	super Pointer `*SuperIndex` // index on the OrderedIndex, IndexBase is the last super index
	final Pointer `*FinalIndex` // index under the OrderedIndex, MultiIndex is the final index

	Root Pointer `*OrderedIndexNode`
	size int
}

func (tree *ByID) init(final *EpIndex) {
	tree.Root = *NewNil()
	tree.final.Set(unsafe.Pointer(final))
	//tree.final = final
	tree.super.Set(unsafe.Pointer(NewSuperIndexByID()))
	//tree.super = NewSuperIndex()
	(*EpIndexBase)(tree.super.Get()).init(final)
	//tree.super.init(final)
}

func (tree *ByID) clear() {
	tree.Clear()
	(*EpIndexBase)(tree.super.Get()).clear()
}

/*generic class*/

const _SizeofSuperIndexByID = unsafe.Sizeof(EpIndexBase{})

func NewSuperIndexByID() *EpIndexBase {
	if alloc == nil {
		return &EpIndexBase{}
	}
	return (*EpIndexBase)(alloc.Allocate(_SizeofSuperIndexByID))
}

/*generic class*/

// OrderedIndexNode is a single element within the tree
type ByIDNode struct {
	Key    uint32
	super  Pointer `*SuperNode`
	final  Pointer `*FinalNode`
	color  colorByID
	Left   Pointer `*OrderedIndexNode`
	Right  Pointer `*OrderedIndexNode`
	Parent Pointer `*OrderedIndexNode`
}

type np_ByID = *ByIDNode

const _SizeofByIDNode = unsafe.Sizeof(ByIDNode{})

func NewByIDNode(key uint32, color colorByID) (n *ByIDNode) {
	if alloc == nil {
		n = new(ByIDNode)
	} else {
		n = np_ByID(alloc.Allocate(_SizeofByIDNode))
	}

	n.Key = key
	n.color = color
	n.super = *NewNil()
	n.final = *NewNil()
	n.Left = *NewNil()
	n.Right = *NewNil()
	n.Parent = *NewNil()

	return n
}

/*generic class*/

/*generic class*/

func (node *ByIDNode) free() {
	if alloc != nil {
		alloc.DeAllocate(unsafe.Pointer(node))
	}
	// else free by golang gc
}

func (node *ByIDNode) value() *item {
	return (*EpIndexBaseNode)(node.super.Get()).value()
	//return node.super.value()
}

type colorByID bool

const (
	blackByID, redByID colorByID = true, false
)

func (tree *ByID) Insert(v item) (IteratorByID, bool) {
	fn, res := (*EpIndex)(tree.final.Get()).insert(v)
	//fn, res := tree.final.insert(v)
	if res {
		return tree.makeIterator(fn), true
	}
	return tree.End(), false
}

func (tree *ByID) insert(v item, fn *EpIndexNode) (*ByIDNode, bool) {
	key := ByIdKeyFunc(v)

	node, res := tree.put(key)
	if !res {
		container.Logger.Warn("#ordered index insert failed")
		return nil, false
	}
	sn, res := (*EpIndexBase)(tree.super.Get()).insert(v, fn)
	//sn, res := tree.super.insert(v, fn)
	if res {
		node.super.Set(unsafe.Pointer(sn))
		//node.super = sn
		node.final.Set(unsafe.Pointer(fn))
		//node.final = fn
		return node, true
	}
	tree.remove(node)
	return nil, false
}

func (tree *ByID) Erase(iter IteratorByID) (itr IteratorByID) {
	itr = iter
	itr.Next()
	(*EpIndex)(tree.final.Get()).erase((*EpIndexNode)(iter.node.final.Get()))
	//tree.final.erase(iter.node.final)
	return
}

func (tree *ByID) Erases(first, last IteratorByID) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *ByID) erase(n *ByIDNode) {
	tree.remove(n)
	(*EpIndexBase)(tree.super.Get()).erase((*EpIndexBaseNode)(n.super.Get()))
	//tree.super.erase(n.super)
	n.super.Set(nil)
	//n.super = nil
	n.final.Set(nil)
	//n.final = nil
}

func (tree *ByID) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(IteratorByID); ok {
		tree.Erase(itr)
	} else {
		(*EpIndexBase)(tree.super.Get()).erase_(iter)
		//tree.super.erase_(iter)
	}
}

func (tree *ByID) Modify(iter IteratorByID, mod func(*item)) bool {
	if _, b := (*EpIndex)(tree.final.Get()).modify(mod, (*EpIndexNode)(iter.node.final.Get())); b {
		//if _, b := tree.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (tree *ByID) modify(n *ByIDNode) (*ByIDNode, bool) {
	n.Key = ByIdKeyFunc(*n.value())

	if !tree.inPlace(n) {
		tree.remove(n)
		node, res := tree.put(n.Key)
		if !res {
			container.Logger.Warn("#ordered index modify failed")
			(*EpIndexBase)(tree.super.Get()).erase((*EpIndexBaseNode)(n.super.Get()))
			//tree.super.erase(n.super)
			return nil, false
		}

		node.super.Forward(&n.super)
		node.final.Forward(&n.final)
		n = node
	}

	if sn, res := (*EpIndexBase)(tree.super.Get()).modify((*EpIndexBaseNode)(n.super.Get())); !res {
		//if sn, res := tree.super.modify(n.super); !res {
		tree.remove(n)
		return nil, false
	} else {
		n.super.Set(unsafe.Pointer(sn))
		//n.super = sn
	}

	return n, true
}

func (tree *ByID) modify_(iter multiindex.IteratorType, mod func(*item)) bool {
	if itr, ok := iter.(IteratorByID); ok {
		return tree.Modify(itr, mod)
	} else {
		return (*EpIndexBase)(tree.super.Get()).modify_(iter, mod)
		//return tree.super.modify_(iter, mod)
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByID) Find(key uint32) IteratorByID {
	if false {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && ByIdCompare(key, lower.Key()) == 0 {
			return lower
		}
		return tree.End()
	} else {
		if node := tree.lookup(key); node != nil {
			return IteratorByID{tree, node, betweenByID}
		}
		return tree.End()
	}
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (tree *ByID) LowerBound(key uint32) IteratorByID {
	result := tree.End()
	node := np_ByID(tree.Root.Get())

	if node == nil {
		return result
	}

	for {
		if ByIdCompare(key, node.Key) > 0 {
			if !node.Right.IsNil() {
				//if node.Right != nil {
				node = np_ByID(node.Right.Get())
				//node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			//result.node = node
			result.position = betweenByID
			if !node.Left.IsNil() {
				//if node.Left != nil {
				node = np_ByID(node.Left.Get())
				//node = node.Left
			} else {
				return result
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (tree *ByID) UpperBound(key uint32) IteratorByID {
	result := tree.End()
	node := np_ByID(tree.Root.Get())

	if node == nil {
		return result
	}

	for {
		if ByIdCompare(key, node.Key) >= 0 {
			if !node.Right.IsNil() {
				//if node.Right != nil {
				node = np_ByID(node.Right.Get())
				//node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			//result.node = node
			result.position = betweenByID
			if !node.Left.IsNil() {
				//if node.Left != nil {
				node = np_ByID(node.Left.Get())
				//node = node.Left
			} else {
				return result
			}
		}
	}
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByID) Remove(key uint32) {
	if false {
		for lower := tree.LowerBound(key); lower.position != endByID; {
			if ByIdCompare(lower.Key(), key) == 0 {
				node := lower.node
				lower.Next()
				tree.remove(node)
			} else {
				break
			}
		}
	} else {
		node := tree.lookup(key)
		tree.remove(node)
	}
}

func (tree *ByID) put(key uint32) (*ByIDNode, bool) {
	var insertedNode *ByIDNode
	if tree.Root.IsNil() {
		//if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		ByIdCompare(key, key)
		tree.Root.Set(unsafe.Pointer(NewByIDNode(key, redByID)))
		//tree.Root = &OrderedIndexNode{Key: key, color: red}
		insertedNode = np_ByID(tree.Root.Get())
	} else {
		node := np_ByID(tree.Root.Get())
		loop := true
		if false {
			for loop {
				compare := ByIdCompare(key, node.Key)
				//compare := Comparator(key, node.Key)
				switch {
				case compare < 0:
					if node.Left.IsNil() {
						//if node.Left == nil {
						node.Left.Set(unsafe.Pointer(NewByIDNode(key, redByID)))
						//node.Left = NewOrderedIndexNode(key, red)
						insertedNode = np_ByID(node.Left.Get())
						//insertedNode = node.Left
						loop = false
					} else {
						node = np_ByID(node.Left.Get())
						//node = node.Left
					}
				case compare >= 0:
					if node.Right.IsNil() {
						//if node.Right == nil {
						node.Right.Set(unsafe.Pointer(NewByIDNode(key, redByID)))
						//node.Right = NewOrderedIndexNode(key, red)
						insertedNode = np_ByID(node.Right.Get())
						//insertedNode = node.Right
						loop = false
					} else {
						node = np_ByID(node.Right.Get())
						//node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := ByIdCompare(key, node.Key)
				//compare := Comparator(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					//node.Key = key
					return node, false
					//return node, false
				case compare < 0:
					if node.Left.IsNil() {
						//if node.Left == nil {
						node.Left.Set(unsafe.Pointer(NewByIDNode(key, redByID)))
						//node.Left = NewOrderedIndexNode(key, red)
						insertedNode = np_ByID(node.Left.Get())
						//insertedNode = node.Left
						loop = false
					} else {
						node = np_ByID(node.Left.Get())
						//node = node.Left
					}
				case compare > 0:
					if node.Right.IsNil() {
						//if node.Right == nil {
						node.Right.Set(unsafe.Pointer(NewByIDNode(key, redByID)))
						//node.Right = NewOrderedIndexNode(key, red)
						insertedNode = np_ByID(node.Right.Get())
						//insertedNode = node.Right
						loop = false
					} else {
						node = np_ByID(node.Right.Get())
						//node = node.Right
					}
				}
			}
		}
		insertedNode.Parent.Set(unsafe.Pointer(node))
		//insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++

	return insertedNode, true
}

func (tree *ByID) swapNode(node *ByIDNode, pred *ByIDNode) {
	if node == pred {
		return
	}

	tmp := ByIDNode{color: pred.color}
	tmp.Left.Forward(&pred.Left)
	tmp.Right.Forward(&pred.Right)
	tmp.Parent.Forward(&pred.Parent)
	//tmp := OrderedIndexNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

	pred.color = node.color
	node.color = tmp.color

	pred.Right.Forward(&node.Right)
	if !pred.Right.IsNil() {
		//if pred.Right != nil {
		np_ByID(pred.Right.Get()).Parent.Set(unsafe.Pointer(pred))
		//pred.Right.Parent = pred
	}
	node.Right.Forward(&tmp.Right)
	if !node.Right.IsNil() {
		//if node.Right != nil {
		np_ByID(pred.Right.Get()).Parent.Set(unsafe.Pointer(node))
		//node.Right.Parent = node
	}

	if np_ByID(pred.Parent.Get()) == node {
		//if pred.Parent == node {
		pred.Left.Set(unsafe.Pointer(node))
		//pred.Left = node
		node.Left.Forward(&tmp.Left)
		//node.Left = tmp.Left
		if !node.Left.IsNil() {
			//if node.Left != nil {
			np_ByID(node.Left.Get()).Parent.Set(unsafe.Pointer(node))
			//node.Left.Parent = node
		}

		pred.Parent.Forward(&node.Parent)
		//pred.Parent = node.Parent
		if !pred.Parent.IsNil() {
			//if pred.Parent != nil {
			if np_ByID(np_ByID(pred.Parent.Get()).Left.Get()) == node {
				//if pred.Parent.Left == node {
				np_ByID(pred.Parent.Get()).Left.Set(unsafe.Pointer(pred))
				//pred.Parent.Left = pred
			} else {
				np_ByID(pred.Parent.Get()).Right.Set(unsafe.Pointer(pred))
				//pred.Parent.Right = pred
			}
		} else {
			tree.Root.Set(unsafe.Pointer(pred))
			//tree.Root = pred
		}
		node.Parent.Set(unsafe.Pointer(pred))
		//node.Parent = pred

	} else {
		pred.Left.Forward(&node.Left)
		if !pred.Left.IsNil() {
			//if pred.Left != nil {
			np_ByID(pred.Left.Get()).Parent.Set(unsafe.Pointer(pred))
			//pred.Left.Parent = pred
		}
		node.Left.Forward(&tmp.Left)
		if !node.Left.IsNil() {
			//if node.Left != nil {
			np_ByID(pred.Left.Get()).Parent.Set(unsafe.Pointer(node))
			//node.Left.Parent = node
		}

		pred.Parent.Forward(&node.Parent)
		if !pred.Parent.IsNil() {
			if np_ByID(np_ByID(pred.Parent.Get()).Left.Get()) == node {
				//if pred.Parent.Left == node {
				np_ByID(pred.Parent.Get()).Left.Set(unsafe.Pointer(pred))
				//pred.Parent.Left = pred
			} else {
				np_ByID(pred.Parent.Get()).Right.Set(unsafe.Pointer(pred))
				//pred.Parent.Right = pred
			}
		} else {
			tree.Root.Set(unsafe.Pointer(pred))
			//tree.Root = pred
		}

		node.Parent.Forward(&tmp.Parent)
		if !node.Parent.IsNil() {
			//if node.Parent != nil {
			if np_ByID(np_ByID(node.Parent.Get()).Left.Get()) == pred {
				//if node.Parent.Left == pred {
				np_ByID(node.Parent.Get()).Left.Set(unsafe.Pointer(node))
				//node.Parent.Left = node
			} else {
				np_ByID(node.Parent.Get()).Right.Set(unsafe.Pointer(node))
				//node.Parent.Right = node
			}
		} else {
			tree.Root.Set(unsafe.Pointer(node))
			//tree.Root = node
		}
	}
}

func (tree *ByID) remove(node *ByIDNode) {
	var child *ByIDNode
	if node == nil {
		return
	}
	if !node.Left.IsNil() && !node.Right.IsNil() {
		//if node.Left != nil && node.Right != nil {
		pred := np_ByID(node.Left.Get()).maximumNode()
		//pred := node.Left.maximumNode()
		tree.swapNode(node, pred)
	}
	if node.Left.IsNil() || node.Right.IsNil() {
		//if node.Left == nil || node.Right == nil {
		if node.Right.IsNil() {
			//if node.Right == nil {
			child = np_ByID(node.Left.Get())
			//child = node.Left
		} else {
			child = np_ByID(node.Right.Get())
			//child = node.Right
		}
		if node.color == blackByID {
			node.color = nodeColorByID(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent.IsNil() && child != nil {
			//if node.Parent == nil && child != nil {
			child.color = blackByID
		}
	}
	tree.size--
	node.free()
}

func (tree *ByID) lookup(key uint32) *ByIDNode {
	node := np_ByID(tree.Root.Get())
	//node := tree.Root
	for node != nil {
		compare := ByIdCompare(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = np_ByID(node.Left.Get())
			//node = node.Left
		case compare > 0:
			node = np_ByID(node.Right.Get())
			//node = node.Right
		}
	}
	return nil
}

// Empty returns true if tree does not contain any nodes
func (tree *ByID) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *ByID) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *ByID) Keys() []uint32 {
	keys := make([]uint32, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *ByID) Values() []item {
	values := make([]item, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *ByID) Left() *ByIDNode {
	var parent *ByIDNode
	current := np_ByID(tree.Root.Get())
	//current := tree.Root
	for current != nil {
		parent = current
		current = np_ByID(current.Left.Get())
		//current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *ByID) Right() *ByIDNode {
	var parent *ByIDNode
	current := np_ByID(tree.Root.Get())
	//current := tree.Root
	for current != nil {
		parent = current
		current = np_ByID(current.Right.Get())
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *ByID) Clear() {
	if alloc != nil {
		//TODO DeAllocator
	}
	tree.Root.Set(nil)
	//tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *ByID) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		outputByID(np_ByID(tree.Root.Get()), "", true, &str)
		//output(tree.Root, "", true, &str)
	}
	return str
}

func (node *ByIDNode) String() string {
	if !node.color {
		return fmt.Sprintf("(%v,%v)", node.Key, "red")
	}
	return fmt.Sprintf("(%v)", node.Key)
}

func outputByID(node *ByIDNode, prefix string, isTail bool, str *string) {
	if !node.Right.IsNil() {
		//if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		outputByID(np_ByID(node.Right.Get()), newPrefix, false, str)
		//output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if !node.Left.IsNil() {
		//if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		outputByID(np_ByID(node.Left.Get()), newPrefix, true, str)
		//output(node.Left, newPrefix, true, str)
	}
}

func (node *ByIDNode) grandparent() *ByIDNode {
	if node != nil && !node.Parent.IsNil() {
		//if node != nil && node.Parent != nil {
		return np_ByID(np_ByID(node.Parent.Get()).Parent.Get())
		//return node.Parent.Parent
	}
	return nil
}

func (node *ByIDNode) uncle() *ByIDNode {
	if node == nil || node.Parent.IsNil() || np_ByID(node.Parent.Get()).Parent.IsNil() {
		//if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return np_ByID(node.Parent.Get()).sibling()
	//return node.Parent.sibling()
}

func (node *ByIDNode) sibling() *ByIDNode {
	if node == nil || node.Parent.IsNil() {
		//if node == nil || node.Parent == nil {
		return nil
	}
	if node == np_ByID(np_ByID(node.Parent.Get()).Left.Get()) {
		//if node == node.Parent.Left {
		return np_ByID(np_ByID(node.Parent.Get()).Right.Get())
		//return node.Parent.Get().Right
	}
	return np_ByID(np_ByID(node.Parent.Get()).Left.Get())
	//return node.Parent.Left
}

func (node *ByIDNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right.IsNil() && node.Left.IsNil() {
		//if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *ByID) rotateLeft(node *ByIDNode) {
	right := np_ByID(node.Right.Get())
	tree.replaceNode(node, right)
	node.Right.Forward(&right.Left)
	if !right.Left.IsNil() {
		//if right.Left != nil {
		np_ByID(right.Left.Get()).Parent.Set(unsafe.Pointer(node))
		//right.Left.Parent = node
	}
	right.Left.Set(unsafe.Pointer(node))
	//right.Left = node
	node.Parent.Set(unsafe.Pointer(right))
	//node.Parent = right
}

func (tree *ByID) rotateRight(node *ByIDNode) {
	left := np_ByID(node.Left.Get())
	//left := node.Left
	tree.replaceNode(node, left)
	node.Left.Forward(&left.Right)
	if !left.Right.IsNil() {
		//if left.Right != nil {
		np_ByID(left.Right.Get()).Parent.Set(unsafe.Pointer(node))
		//left.Right.Parent = node
	}
	left.Right.Set(unsafe.Pointer(node))
	//left.Right = node
	node.Parent.Set(unsafe.Pointer(left))
}

func (tree *ByID) replaceNode(old *ByIDNode, new *ByIDNode) {
	if old.Parent.IsNil() {
		//if old.Parent == nil {
		tree.Root.Set(unsafe.Pointer(new))
		//tree.Root = new
	} else {
		if old == np_ByID(np_ByID(old.Parent.Get()).Left.Get()) {
			//if old == old.Parent.Left {
			np_ByID(old.Parent.Get()).Left.Set(unsafe.Pointer(new))
			//old.Parent.Left = new
		} else {
			np_ByID(old.Parent.Get()).Right.Set(unsafe.Pointer(new))
			//old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent.Forward(&old.Parent)
	}
}

func (tree *ByID) insertCase1(node *ByIDNode) {
	if node.Parent.IsNil() {
		//if node.Parent == nil {
		node.color = blackByID
	} else {
		tree.insertCase2(node)
	}
}

func (tree *ByID) insertCase2(node *ByIDNode) {
	if nodeColorByID(np_ByID(node.Parent.Get())) == blackByID {
		//if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *ByID) insertCase3(node *ByIDNode) {
	uncle := node.uncle()
	if nodeColorByID(uncle) == redByID {
		np_ByID(node.Parent.Get()).color = blackByID
		//node.Parent.color = black
		uncle.color = blackByID
		node.grandparent().color = redByID
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *ByID) insertCase4(node *ByIDNode) {
	grandparent := node.grandparent()
	if node == np_ByID(np_ByID(node.Parent.Get()).Right.Get()) && node.Parent.Get() == grandparent.Left.Get() {
		//if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(np_ByID(node.Parent.Get()))
		//tree.rotateLeft(node.Parent)
		node = np_ByID(node.Left.Get())
		//node = node.Left
	} else if node == np_ByID(np_ByID(node.Parent.Get()).Left.Get()) && node.Parent.Get() == grandparent.Right.Get() {
		//} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(np_ByID(node.Parent.Get()))
		//tree.rotateRight(node.Parent)
		node = np_ByID(node.Right.Get())
		//node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *ByID) insertCase5(node *ByIDNode) {
	np_ByID(node.Parent.Get()).color = blackByID
	//node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = redByID
	if node == np_ByID(np_ByID(node.Parent.Get()).Left.Get()) && node.Parent.Get() == grandparent.Left.Get() {
		//if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == np_ByID(np_ByID(node.Parent.Get()).Right.Get()) && node.Parent.Get() == grandparent.Right.Get() {
		//} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *ByIDNode) maximumNode() *ByIDNode {
	if node == nil {
		return nil
	}
	for !node.Right.IsNil() {
		//for node.Right != nil {
		node = np_ByID(node.Right.Get())
		//node = node.Right
	}
	return node
}

func (tree *ByID) deleteCase1(node *ByIDNode) {
	if node.Parent.IsNil() {
		//if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *ByID) deleteCase2(node *ByIDNode) {
	sibling := node.sibling()
	if nodeColorByID(sibling) == redByID {
		np_ByID(node.Parent.Get()).color = redByID
		//node.Parent.color = red
		sibling.color = blackByID
		if node == np_ByID(np_ByID(node.Parent.Get()).Left.Get()) {
			//if node == node.Parent.Left {
			tree.rotateLeft(np_ByID(node.Parent.Get()))
			//tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(np_ByID(node.Parent.Get()))
			//tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *ByID) deleteCase3(node *ByIDNode) {
	sibling := node.sibling()
	if nodeColorByID(np_ByID(node.Parent.Get())) == blackByID &&
		//if nodeColor(node.Parent) == black &&
		nodeColorByID(sibling) == blackByID &&
		//nodeColor(sibling) == black &&
		nodeColorByID(np_ByID(sibling.Left.Get())) == blackByID &&
		//nodeColor(sibling.Left) == black &&
		nodeColorByID(np_ByID(sibling.Right.Get())) == blackByID {
		//nodeColor(sibling.Right) == black {
		sibling.color = redByID
		tree.deleteCase1(np_ByID(node.Parent.Get()))
		//tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *ByID) deleteCase4(node *ByIDNode) {
	sibling := node.sibling()
	if nodeColorByID(np_ByID(node.Parent.Get())) == redByID &&
		//if nodeColor(node.Parent) == red &&
		nodeColorByID(sibling) == blackByID &&
		nodeColorByID(np_ByID(sibling.Left.Get())) == blackByID &&
		//nodeColor(sibling.Left) == black &&
		nodeColorByID(np_ByID(sibling.Right.Get())) == blackByID {
		//nodeColor(sibling.Right) == black {
		sibling.color = redByID
		np_ByID(node.Parent.Get()).color = blackByID
		//node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *ByID) deleteCase5(node *ByIDNode) {
	sibling := node.sibling()
	if node == np_ByID(np_ByID(node.Parent.Get()).Left.Get()) &&
		//if node == node.Parent.Left &&
		nodeColorByID(sibling) == blackByID &&
		nodeColorByID(np_ByID(sibling.Left.Get())) == redByID &&
		//nodeColor(sibling.Left) == red &&
		nodeColorByID(np_ByID(sibling.Right.Get())) == blackByID {
		//nodeColor(sibling.Right) == black {
		sibling.color = redByID
		np_ByID(sibling.Left.Get()).color = blackByID
		//sibling.Left.color = black
		tree.rotateRight(sibling)
	} else if node == np_ByID(np_ByID(node.Parent.Get()).Right.Get()) &&
		//} else if node == node.Parent.Right &&
		nodeColorByID(sibling) == blackByID &&
		nodeColorByID(np_ByID(sibling.Right.Get())) == redByID &&
		//nodeColor(sibling.Right) == red &&
		nodeColorByID(np_ByID(sibling.Left.Get())) == blackByID {
		//nodeColor(sibling.Left) == black {
		sibling.color = redByID
		np_ByID(sibling.Right.Get()).color = blackByID
		//sibling.Right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *ByID) deleteCase6(node *ByIDNode) {
	sibling := node.sibling()
	sibling.color = nodeColorByID(np_ByID(node.Parent.Get()))
	//sibling.color = nodeColor(node.Parent)
	np_ByID(node.Parent.Get()).color = blackByID
	//node.Parent.color = black
	if node == np_ByID(np_ByID(node.Parent.Get()).Left.Get()) && nodeColorByID(np_ByID(sibling.Right.Get())) == redByID {
		//if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		np_ByID(sibling.Right.Get()).color = blackByID
		//sibling.Right.color = black
		tree.rotateLeft(np_ByID(node.Parent.Get()))
		//tree.rotateLeft(node.Parent)
	} else if nodeColorByID(np_ByID(sibling.Left.Get())) == redByID {
		//} else if nodeColor(sibling.Left) == red {
		np_ByID(sibling.Left.Get()).color = blackByID
		//sibling.Left.color = black
		tree.rotateRight(np_ByID(node.Parent.Get()))
		//tree.rotateRight(node.Parent)
	}
}

func nodeColorByID(node *ByIDNode) colorByID {
	if node == nil {
		return blackByID
	}
	return node.color
}

//////////////iterator////////////////

func (tree *ByID) makeIterator(fn *EpIndexNode) IteratorByID {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(np_ByID); ok {
			return IteratorByID{tree: tree, node: n, position: betweenByID}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

// Iterator holding the iterator's state
type IteratorByID struct {
	tree     *ByID
	node     *ByIDNode
	position positionByID
}

type positionByID byte

const (
	beginByID, betweenByID, endByID positionByID = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *ByID) Iterator() IteratorByID {
	return IteratorByID{tree: tree, node: nil, position: beginByID}
}

func (tree *ByID) Begin() IteratorByID {
	itr := IteratorByID{tree: tree, node: nil, position: beginByID}
	itr.Next()
	return itr
}

func (tree *ByID) End() IteratorByID {
	return IteratorByID{tree: tree, node: nil, position: endByID}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *IteratorByID) Next() bool {
	if iterator.position == endByID {
		goto end
	}
	if iterator.position == beginByID {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if !iterator.node.Right.IsNil() {
		//if iterator.node.Right != nil {
		iterator.node = np_ByID(iterator.node.Right.Get())
		//iterator.node = iterator.node.Right
		for !iterator.node.Left.IsNil() {
			//for iterator.node.Left != nil {
			iterator.node = np_ByID(iterator.node.Left.Get())
			//iterator.node = iterator.node.Left
		}
		goto between
	}
	if !iterator.node.Parent.IsNil() {
		//if iterator.node.Parent != nil {
		node := iterator.node
		for !iterator.node.Parent.IsNil() {
			//for iterator.node.Parent != nil {
			iterator.node = np_ByID(iterator.node.Parent.Get())
			//iterator.node = iterator.node.Parent
			if node == np_ByID(iterator.node.Left.Get()) {
				//if node == iterator.node.Left {
				goto between
			}
			node = iterator.node
		}
	}

end:
	iterator.node = nil
	iterator.position = endByID
	return false

between:
	iterator.position = betweenByID
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *IteratorByID) Prev() bool {
	if iterator.position == beginByID {
		goto begin
	}
	if iterator.position == endByID {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if !iterator.node.Left.IsNil() {
		//if iterator.node.Left != nil {
		iterator.node = np_ByID(iterator.node.Left.Get())
		//iterator.node = iterator.node.Left
		for !iterator.node.Right.IsNil() {
			//for iterator.node.Right != nil {
			iterator.node = np_ByID(iterator.node.Right.Get())
			//iterator.node = iterator.node.Right
		}
		goto between
	}
	if !iterator.node.Parent.IsNil() {
		//if iterator.node.Parent != nil {
		node := iterator.node
		for !iterator.node.Parent.IsNil() {
			//for iterator.node.Parent != nil {
			iterator.node = np_ByID(iterator.node.Parent.Get())
			//iterator.node = iterator.node.Parent
			if node == np_ByID(iterator.node.Right.Get()) {
				//if node == iterator.node.Right {
				goto between
			}
			node = iterator.node
		}
	}

begin:
	iterator.node = nil
	iterator.position = beginByID
	return false

between:
	iterator.position = betweenByID
	return true
}

func (iterator IteratorByID) HasNext() bool {
	return iterator.position != endByID
}

func (iterator *IteratorByID) HasPrev() bool {
	return iterator.position != beginByID
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorByID) Value() item {
	return *iterator.node.value()
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator IteratorByID) Key() uint32 {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *IteratorByID) Begin() {
	iterator.node = nil
	iterator.position = beginByID
}

func (iterator IteratorByID) IsBegin() bool {
	return iterator.position == beginByID
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *IteratorByID) End() {
	iterator.node = nil
	iterator.position = endByID
}

func (iterator IteratorByID) IsEnd() bool {
	return iterator.position == endByID
}

// Delete remove the node which pointed by the iterator
// Modifies the state of the iterator.
func (iterator *IteratorByID) Delete() {
	node := iterator.node
	//iterator.Prev()
	iterator.tree.remove(node)
}

func (tree *ByID) inPlace(n *ByIDNode) bool {
	prev := IteratorByID{tree, n, betweenByID}
	next := IteratorByID{tree, n, betweenByID}
	prev.Prev()
	next.Next()

	var (
		prevResult int
		nextResult int
	)

	if prev.IsBegin() {
		prevResult = 1
	} else {
		prevResult = ByIdCompare(n.Key, prev.Key())
	}

	if next.IsEnd() {
		nextResult = -1
	} else {
		nextResult = ByIdCompare(n.Key, next.Key())
	}

	return (false && prevResult >= 0 && nextResult <= 0) ||
		(!false && prevResult > 0 && nextResult < 0)
}

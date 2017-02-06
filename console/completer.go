package console

// Completer is an interface that can be implemented to allow users custom
// tab-completion rules inside an ugCLI console.
type Completer interface {
	// Complete takes a current command line, and returns a common prefix to
	// expand the command line to, along with a list of potential recommendations.
	// If there is only a single recommendation, then the console will simply
	// set the current line to that recommendation. Otherwise, it will display
	// all possible recommendations below the current line.
	Complete(input string) (prefix string, recommendations []string)
}

// listCompleter is a default, builtin implementation of the completer
// interface. It uses a prefix tree (trie) and recommends all commands
// stored within its list that share a prefix with the given input.
type listCompleter struct {

	// head is the root of the prefix tree.
	head *lcNode

	// list holds all possible commands that the console will tab-complete to.
	list []string
}

// NewListCompleter will take a list of potential command words and produce
// a completer that return all elements of the list that have the given
// line as a prefix.
func NewListCompleter(list []string) Completer {

	// Initialize the head of the prefix tree.
	head := &lcNode{
		children: map[rune]*lcNode{},
		end:      false,
	}

	// Add every word in the list to the prefix tree.
	for _, word := range list {

		// N will be the node we look at, as we descend down the tree.
		n := head

		// For every character in the word:
		for _, char := range word {

			if child, exists := n.children[char]; exists {
				// If this node already has a child with this character, move to the
				// child.
				n = child
			} else {
				// Otherwise, insert a new node here and visit it.
				n.children[char] = &lcNode{
					children: map[rune]*lcNode{},
					end:      false,
				}
				n = n.children[char]
			}
		}

		// At the end of the word, mark the last node as a terminal node for
		// some word.
		n.end = true
	}

	// Store the prefix tree in the completer, along with the word list.
	return &listCompleter{
		head: head,
		list: list,
	}
}

// Complete implements the Completer interface.
func (lc *listCompleter) Complete(input string) (string, []string) {

	// Our node as we traverse the tree.
	n := lc.head
	var exists bool

	// If the given word isn't in the prefix tree at all, return
	// empty results.
	for _, char := range input {
		if n, exists = n.children[char]; !exists {
			return "", []string{}
		}
	}

	// Otherwise, find all matches that are children of the node at the
	// end of the given string (that have the given string as a prefix).
	return lcNodeDFS(n, input, []string{})
}

// lcNodeDFS is a utility function used to traverse the prefix tree.
// Given a current node, the running string we are at, and all previously
// added substrings, it will return all strings stored as children of the
// current node.
func lcNodeDFS(n *lcNode, running string, prev []string) (string, []string) {

	prefix := running

	// If this is a terminal state, add it to the output.
	if n.end {
		prev = append(prev, running)
	}

	// For every child, recursively search for any more strings.
	for c, child := range n.children {
		prefix, prev = lcNodeDFS(child, running+string(c), prev)
	}

	if len(n.children) == 1 && !n.end {
		// If there is only one child, and this is a nonterminal state, then
		// the common prefix is the common prefix of the subtree.
		return prefix, prev
	} else {
		// Otherwise, the common prefix does not contain any children of this node.
		return running, prev
	}
}

// lcNode represents a node used in the prefix tree for a list completer.
type lcNode struct {

	// children stores all child nodes.
	children map[rune]*lcNode

	// end indicates whether the string represented by the path from root to
	// this node is an actual element of the list, rather than just a prefix.
	end bool
}

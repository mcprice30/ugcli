// ugCLI is a library built atop termbox for creating CLI applications.
package console

// The Completer interface will be used in the future
// to allow tab-completion in an ugCLI console.
type Completer interface {
	Complete(input string) (prefix string, recommendations []string)
}

type listCompleter struct {
	head *lcNode
	list []string
}

// A built in Completer implementation that will complete based on matches
// in a list of provided words.
func NewListCompleter(list []string) Completer {
	head := &lcNode{
		children: map[rune]*lcNode{},
		end:      false,
	}

	for _, word := range list {

		n := head

		for _, char := range word {
			if child, exists := n.children[char]; exists {
				n = child
			} else {
				n.children[char] = &lcNode{
					children: map[rune]*lcNode{},
					end:      false,
				}
				n = n.children[char]
			}
		}

		n.end = true
	}

	return &listCompleter{
		head: head,
		list: list,
	}
}

func (lc *listCompleter) Complete(input string) (string, []string) {
	n := lc.head
	var exists bool

	for _, char := range input {
		if n, exists = n.children[char]; !exists {
			return "", []string{}
		}
	}

	return lcNodeDFS(n, input, []string{})
}

func lcNodeDFS(n *lcNode, running string, prev []string) (string, []string) {

	prefix := running

	if n.end {
		prev = append(prev, running)
	}

	for c, child := range n.children {
		prefix, prev = lcNodeDFS(child, running+string(c), prev)
	}

	if len(n.children) == 1 && !n.end {
		return prefix, prev
	} else {
		return running, prev
	}
}

type lcNode struct {
	children map[rune]*lcNode
	end      bool
}

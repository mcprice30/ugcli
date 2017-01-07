// ugCLI is a library built atop termbox for creating CLI applications.
package ugcli

// The Completer interface will be used in the future
// to allow tab-completion in an ugCLI console.

type Completer interface {
	Complete(input string) (string, []string)
}

type ListCompleter struct {
	head *lcNode
	list []string
}

func CreateListCompleter(list []string) Completer {
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

	return &ListCompleter{
		head: head,
		list: list,
	}
}

func (lc *ListCompleter) Complete(input string) (string, []string) {
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

func (n *lcNode) String() string {
	out := ""
	if n.end {
		out = "true"
	} else {
		out = "false"
	}

	for c, child := range n.children {
		out += "\n" + string(c) + ":\n" + child.String()
	}

	out += "\n\n"
	return out
}

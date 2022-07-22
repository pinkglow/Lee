package lee

import (
	"strings"
)

// 前缀树的节点
type node struct {
	pattern  string  // 叶子结点不为空，存储完整路径，其他节点为 nil
	part     string  // 当前节点，存储的 part
	children []*node // 该节点的子节点
	isWild   bool    // 路径包含有 ：或者 * 为 true，否则为 false
}

// matchChild 获取第一个与 part 匹配的路径
func (n *node) matchChild(part string) *node {
	// 从该节点的子节点遍历
	for _, child := range n.children {
		if child.part == part || n.isWild {
			return child
		}
	}
	return nil
}

// 经过解析后得到 parts ["author", "lookcos"]，因此会插入 2 个 node
func (n *node) _insertNode(pattern string, parts []string, height int) {
	// 说明已经插入完毕
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	// 找到当前要插入的那个 part
	part := parts[height]
	// 根据 part 查找它的对应 child
	child := n.matchChild(part)
	// 如果child 为空, 说明需要新建此节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == '*' || part[0] == ':'}
		n.children = append(n.children, child)
	}

	// 接着插入下一个 part
	child._insertNode(pattern, parts, height+1)
}

func (n *node) insertNode(pattern string) {
	parts := parsePattern(pattern)
	n._insertNode(pattern, parts, 0)
}

// HasPrefix 用于判断是否含有某个前缀
func (n *node) HasPrefix(prefix string) bool {
	return n._HasPrefix(prefix)
}

func (n *node) _HasPrefix(prefix string) bool {
	parts := parsePattern(prefix)
	if n._findNode(parts, 0, true) != nil {
		return true
	}
	return false
}

// matchChildren 用来遍历节点 n 下面所有与 part 相匹配的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if part == child.part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// findNode 获取路径对应的 node
func (n *node) _findNode(parts []string, height int, isFindPrefix bool) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		// 说明查找失败
		if n.pattern == "" && isFindPrefix == false {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	// 在与当前节点匹配的所有子节点中搜索
	for _, child := range children {
		result := child._findNode(parts, height+1, isFindPrefix)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) findNode(pattern string) *node {
	parts := parsePattern(pattern)
	return n._findNode(parts, 0, false)
}

// getFullPathNodes 遍历整棵树，找到所有有完整路径的节点
func (n *node) getFullPathNodes(nodes *[]*node) {
	if n.pattern != "" {
		*nodes = append(*nodes, n)
	}
	for _, child := range n.children {
		child.getFullPathNodes(nodes)
	}
}

// parsePattern 解析 pattern 为 parts
func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			// 路径含有通配符，后面的就不用管了
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

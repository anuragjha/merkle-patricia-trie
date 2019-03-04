package p1

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/crypto/sha3"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db   map[string]Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {

	fmt.Println("Geting value of Key: ", key)
	pathLeft := StringToHexArray(key)
	currentNode := mpt.db[mpt.root]

	value, err := mpt.getHelper(currentNode, pathLeft)

	if err != nil {
		fmt.Println("Value:", err)
		return "", err
	} else {
		fmt.Println("Value:", value)
		return value, nil
	}
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	fmt.Println("Key - Value being inserted:", key, " - ", new_value)

	pathLeft := StringToHexArray(key)
	if mpt.root == "" {
		// pathLeft = append(pathLeft, 16)
		// leafNode := Node{}
		// leafNode.node_type = 2
		// leafNode.flag_value.encoded_prefix = compact_encode(pathLeft)

		// leafNode.flag_value.value = new_value
		// mpt.root = mpt.insertHelper(leafNode, nil, "", nil)
		//>>
		mpt.root = mpt.insertHelper1(Node{}, pathLeft, new_value)

	} else {

		// rootNodeHash := mpt.root
		// mpt.root = mpt.insertHelper(mpt.db[rootNodeHash], pathLeft, new_value, []uint8{})
		//>>
		mpt.root = mpt.insertHelper1(mpt.db[mpt.root], pathLeft, new_value)
	}

}

//Delete func
func (mpt *MerklePatriciaTrie) Delete(key string) string {
	pathLeft := StringToHexArray(key)
	currentNode := mpt.db[mpt.root]

	value, err := mpt.delHelper(mpt.root, currentNode, pathLeft, []string{})
	fmt.Println("Key - Value deleted : ", key, " - ", value)
	if err != nil {
		return "path_not_found"
	}
	return value
}

func compact_encode(hex_array []uint8) []uint8 {
	var term int
	if hex_array[len(hex_array)-1] == 16 { //checking if last element in array is 16
		term = 1 // if last element is 16, term = 1 and remove the last element from hex_array
		hex_array = hex_array[:len(hex_array)-1]
	}

	oddlen := len(hex_array) % 2               //checking if length is odd (oddlen = 1) otherwise (oddlen = 0)
	flags := []uint8{(uint8(2*term + oddlen))} //calculating flags value
	//changing hex_array based on value of oddlen
	if oddlen == 1 { //odd       // prefix -> flags so... ( either, 1 - Odd length extension
		hex_array = append(flags, hex_array...) //      3 - Odd length leaf        )
	} else { //when oddlen = 0 // even
		flags = append(flags, uint8(0))
		hex_array = append(flags, hex_array...) //prefix -> flags + 0 so...  (either, 00 -  Even Length Extention)
	} // or 20 - Even length Leaf
	encoded_arr := []uint8{} //array to return
	for i := 0; i < len(hex_array); i += 2 {
		encoded_arr = append(encoded_arr, 16*hex_array[i]+hex_array[i+1])
	}

	return encoded_arr
} // closing compact_encode

func compact_decode(encoded_arr []uint8) []uint8 {
	decoded_arr := AsciiArrayToHexArray(encoded_arr) //converting ascii array to hex array

	prefix := decoded_arr[0]
	switch prefix {
	case 0:
		decoded_arr = decoded_arr[2:]
	case 1:
		decoded_arr = decoded_arr[1:]
	case 2:
		decoded_arr = decoded_arr[2:]
	case 3:
		decoded_arr = decoded_arr[1:]
	}
	return decoded_arr
} //closing compact_decode

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value + string(node.flag_value.encoded_prefix)
	}
	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {

	mpt.db = make(map[string]Node)
	mpt.root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func is_leaf_node(encoded_arr []uint8) bool {
	return (encoded_arr[0]/16 == 2) || (encoded_arr[0]/16 == 3)
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}

// StringToHexArray takes in string as input and returns a splitted-hex(0-15) array : 1st
func StringToHexArray(s string) []uint8 {
	uint8Array := []uint8(s)
	hexArray := AsciiArrayToHexArray(uint8Array)
	return hexArray
}

// AsciiArrayToHexArray takes in string as input and returns a splitted-hex(0-15) array : 1st
func AsciiArrayToHexArray(encoded_arr []uint8) []uint8 {
	hexArray := []uint8{}
	for _, element := range encoded_arr {
		div, mod := element/16, element%16
		hexArray = append(hexArray, div)
		hexArray = append(hexArray, mod)
	}
	return hexArray
}

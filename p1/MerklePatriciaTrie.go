package p1

import (
	"encoding/hex"
	"errors"
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

	fmt.Println("VALUE to be searched:", key)
	// return "", errors.New("path_not_found")
	fmt.Println("IN GET:", key)
	pathLeft := StringToHexArray(key)
	fmt.Println("input Hex:", pathLeft)
	currentNode := mpt.db[mpt.root]
	fmt.Printf("root node : %+v", currentNode)
	finalValue, err := mpt.getHelper(currentNode, pathLeft)

	//if finalNode.
	fmt.Println("finalValue:", finalValue)

	if err != nil {
		return "", err
	} else {
		return finalValue, nil
	}

}

////////getHELPER

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	fmt.Println("KEY being inserted:", key)
	fmt.Println("VALUE being inserted:", new_value)

	//hex_key_array := encode(toByteArray(key))
	hex_key_array := StringToHexArray(key)
	if mpt.root == "" {
		fmt.Println("No root")
		hex_key_array = append(hex_key_array, 16)
		nodeL := Node{}
		nodeL.node_type = 2
		nodeL.flag_value.encoded_prefix = compact_encode(hex_key_array)
		fmt.Println("Encoded prefix :", nodeL.flag_value.encoded_prefix)
		nodeL.flag_value.value = new_value
		mpt.root = mpt.InsertRecursive(nodeL, nil, "", nil)
		fmt.Println("Inserted in root")
	} else {
		root_node_hash := mpt.root
		fmt.Println("Root hash before", mpt.root)
		mpt.root = mpt.InsertRecursive(mpt.db[root_node_hash], hex_key_array, new_value, []uint8{})
		fmt.Println("Root updated")
		fmt.Println("Root hash after", mpt.root)
	}

}

/// insert Helper

//createLeafNode
func createLeafNode(path []uint8, new_value string) (string, Node) {
	flagStruct := Flag_value{path, new_value}
	n1 := &Node{
		node_type:  2,
		flag_value: flagStruct,
	}
	nHash := n1.hash_node()
	return nHash, *n1
}

//createExtensionNode
func createExtensionNode(path []uint8, next_hash_value string) (string, Node) {
	flagStruct := Flag_value{path, next_hash_value}
	n1 := &Node{
		node_type:  2,
		flag_value: flagStruct,
	}
	nHash := n1.hash_node()
	return nHash, *n1
}

//createBranchNode
func createBranchNode(path []uint8, new_value string) (string, Node) {
	n1 := &Node{
		node_type:    1,
		branch_value: [17]string{},
	}
	nHash := n1.hash_node()
	return nHash, *n1
}

//Delete func
func (mpt *MerklePatriciaTrie) Delete(key string) string {
	fmt.Println("VALUE to be Deleted :", key)
	pathLeft := StringToHexArray(key)
	fmt.Println("in Del - input Hex:", pathLeft)
	currentNode := mpt.db[mpt.root]
	fmt.Printf("in Del - root node : %+v", currentNode)
	//sendinf root hash in hashstack
	value, err := mpt.delHelper(mpt.root, currentNode, pathLeft, []string{})
	fmt.Println("in Del -deleted value: ", value)
	if err != nil {
		return value
	}
	return ""

}

// deleteHelper & rearrange

func compact_encode(hex_array []uint8) []uint8 {
	var term int
	if hex_array[len(hex_array)-1] == 16 { //checking if last element in array is 16
		term = 1 // if last element is 16, term = 1 and remove the last element from hex_array
		hex_array = hex_array[:len(hex_array)-1]
		if len(hex_array) == 0 {
			fmt.Println("Empty leaf in compact encode")
		}
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
	//fmt.Println("ascii array :", ascii_array)
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
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
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

func main() { //before pepsi
	mpt := &MerklePatriciaTrie{make(map[string]Node), ""}
	//mpt.CreateTestMpt()

	//fmt.Printf("%+v\n", mpt)
	// fmt.Printf("mpt-blank: %+v\n", mpt)
	// ////insert do -> verb
	// mpt.Insert("do", "verb")
	// fmt.Printf("mpt-do: %+v\n", mpt)
	////////////

	fmt.Println("inserting do")
	mpt.Insert("do", "verb")
	fmt.Println("inserting dog")
	mpt.Insert("dog", "puppy")
	fmt.Println("inserting doge")
	mpt.Insert("doge", "coin")
	fmt.Println("inserting horse")
	mpt.Insert("horse", "stallion")

	fmt.Println("####################################")
	fmt.Println("====================================")
	v, e := mpt.Get("dog")
	if e != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("value1 : ", v)
	}

	fmt.Println("====================================")
	fmt.Println("====================================")

	value := mpt.Delete("dog")
	fmt.Println("value2 : ", value)

	fmt.Println("====================================")
	fmt.Println("====================================")

	v1, e1 := mpt.Get("do")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("do : ", v1)
	}

	fmt.Println("====================================")
	fmt.Println("====================================")
	fmt.Println("====================================")

	fmt.Println("inserting dog again.............................")
	mpt.Insert("dog", "puppy")

	fmt.Println("====================================")
	fmt.Println("====================================")
	fmt.Println("====================================")

	// fmt.Println("====================================")
	// fmt.Println("====================================")
	// fmt.Println("====================================")

	v1, e1 = mpt.Get("dog")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("dog : ", v1)
	}

	fmt.Println("====================================")
	// fmt.Println("====================================")
	// fmt.Println("====================================")
	// fmt.Println("====================================")
	v2, e2 := mpt.Get("do")
	if e2 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("do : ", v2)
	}
	fmt.Println("====================================")
	// fmt.Println("====================================")
	// v1, e1 = mpt.Get("doge")
	// if e1 != nil {
	// 	fmt.Println("error in GET method")
	// } else {
	// 	fmt.Println("doge : ", v1)
	// }

	// v1, e1 = mpt.Get("horse")

	// if e1 != nil {
	// 	fmt.Println("error in GET method")
	// } else {
	// 	fmt.Println("horse : ", v1)
	// }

	v1, e1 = mpt.Get("doge")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("doge : ", v1)
	}

	fmt.Println("====================================")
	v1, e1 = mpt.Get("horse")
	if e1 != nil {
		fmt.Println("error in GET method")
	} else {
		fmt.Println("horse : ", v1)
	}
	fmt.Println("====================================")
	fmt.Println("====================================")
	// fmt.Println(mpt.String())
	// v1, e1 = mpt.Get("horse")
	// if e1 != nil {
	// 	fmt.Println("error in GET method")
	// } else {
	// 	fmt.Println("horse : ", v1)
	// }
	// fmt.Println("====================================")
	// fmt.Println("====================================")
	// fmt.Println("====================================")
	// fmt.Println("====================================")

	// fmt.Println("inserting horse")
	// mpt.Insert("horse", "stallion")

	// fmt.Println("inserting dorg")
	// mpt.Insert("dorg", "purppy")

	// vd, ed := mpt.Get("dorg")
	// if ed != nil {
	// 	fmt.Println("error in GET method")
	// } else {
	// 	fmt.Println("dorg : ", vd)
	// }

	// // fmt.Println("inserting dog")
	// // mpt.Insert("dog", "puppy")

	// vd, ed = mpt.Get("dog")
	// if ed != nil {
	// 	fmt.Println("error in GET method")
	// } else {
	// 	fmt.Println("dog : ", vd)
	// }

	// fmt.Println("EOP")
	// fmt.Println("====================================")
	// fmt.Println("####################################")

}

//main to insert node
// func main() {
// 	////initialize mpt to nil
// 	mpt := &MerklePatriciaTrie{make(map[string]Node), ""}
// 	fmt.Printf("mpt-blank: %+v\n", mpt)
// 	////insert do -> verb
// 	mpt.Insert("do", "verb")
// 	fmt.Printf("mpt-do: %+v\n", mpt)
// 	////insert dog -> puppy
// 	mpt.Insert("dog", "puppy")
// 	fmt.Printf("mpt-dog: %+v\n", mpt)

// }

//main to test encode decode
// func main() {
// 	bla := StringToHexArray("dog")
// 	//bla = append(bla, uint8(16))
// 	fmt.Println("bla :", bla)

// 	bla1 := compact_encode(bla)
// 	fmt.Println("bla1 :", bla1)

// 	bla2 := compact_decode(bla1)
// 	fmt.Println("bla2 :", bla2)

// abla := []uint8{1, 6, 1}
//abla := []uint8{32, 100, 111}
//abla := []uint8{6, 4, 6, 15, 16}
//abla1 := compact_encode(abla)
//fmt.Println("abla1=ascii :", abla1)
// abla2 := compact_decode(abla1)
// fmt.Println("abla2=hex : ", abla2)

// 	//root node
// 	// flagStruct := Flag_value{bla1, "verb"}
// 	// n1 := Node{
// 	// 	node_type:  2,
// 	// 	flag_value: flagStruct,
// 	// }
// 	// hashVal := n1.hash_node() //hash to put in db !!!!
// 	// fmt.Println("hashVal", hashVal)

// }

func (mpt *MerklePatriciaTrie) CreateTestMpt() error {
	mpt.db = make(map[string]Node)

	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
		encoded_prefix: compact_encode([]uint8{5, 16}),
		value:          "coin",
	}

	nodeG := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeG,
	}
	hashNodeG := nodeG.hash_node()
	mpt.db[hashNodeG] = nodeG

	fmt.Println("Compact_encode([]uint8{5, 16}:", compact_encode([]uint8{5, 16}))

	flagValueNodeF := Flag_value{
		encoded_prefix: nil,
		value:          "puppy",
	}

	nodeF := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeF,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeG, "", "", "", "", "", "", "", "", "", "puppy"},
	}
	hashNodeF := nodeF.hash_node()
	mpt.db[hashNodeF] = nodeF

	flagValueNodeE := Flag_value{
		encoded_prefix: compact_encode([]uint8{7}),
		value:          hashNodeF,
	}

	nodeE := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.hash_node()
	mpt.db[hashNodeE] = nodeE

	flagValueNodeD := Flag_value{
		encoded_prefix: nil,
		value:          "verb",
	}

	nodeD := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeD,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeE, "", "", "", "", "", "", "", "", "", "verb"},
	}
	hashNodeD := nodeD.hash_node()
	mpt.db[hashNodeD] = nodeD

	flagValueNodeB := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15}),
		value:          hashNodeD,
	}

	nodeB := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.hash_node()
	mpt.db[hashNodeB] = nodeB

	flagValueNodeC := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15, 7, 2, 7, 3, 6, 5, 16}),
		value:          "", //"stallion",
	}

	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.hash_node()
	mpt.db[hashNodeC] = nodeC

	flagValueNodeA := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeA := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeA,
		branch_value: [17]string{"", "", "", "", hashNodeB, "", "", "", hashNodeC, "", "", "", "", "", "", "", ""},
	}
	hashNodeA := nodeA.hash_node()
	mpt.db[hashNodeA] = nodeA

	flagValueRoot := Flag_value{
		encoded_prefix: compact_encode([]uint8{6}),
		value:          hashNodeA,
	}

	nodeRoot := Node{
		node_type:  2, //Extension Root
		flag_value: flagValueRoot,
	}
	hashNodeRoot := nodeRoot.hash_node()
	mpt.db[hashNodeRoot] = nodeRoot
	mpt.root = hashNodeRoot

	////////////////////////////////////////

	//Add another node to MPT

	//mpt.root = Compact_encode(key) //a
	return errors.New("Problem occured while creating Root Node")
}

func (mpt *MerklePatriciaTrie) CreateTestMpt2() error {
	mpt.db = make(map[string]Node)

	// 0: Null, 1: Branch, 2: Ext or Leaf
	flagValueNodeG := Flag_value{
		encoded_prefix: compact_encode([]uint8{5, 16}),
		value:          "coin",
	}

	nodeG := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeG,
	}
	hashNodeG := nodeG.hash_node()
	mpt.db[hashNodeG] = nodeG

	fmt.Println("Compact_encode([]uint8{5, 16}:", compact_encode([]uint8{5, 16}))

	flagValueNodeF := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeF := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeF,
		branch_value: [17]string{"", "", "", "", "", "", hashNodeG, "", "", "", "", "", "", "", "", "", "puppy"},
	}
	hashNodeF := nodeF.hash_node()
	mpt.db[hashNodeF] = nodeF

	flagValueNodeE := Flag_value{
		encoded_prefix: compact_encode([]uint8{7}),
		value:          hashNodeF,
	}

	nodeE := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeE,
	}
	hashNodeE := nodeE.hash_node()
	mpt.db[hashNodeE] = nodeE

	flagValueNodeJ := Flag_value{
		encoded_prefix: compact_encode([]uint8{8}),
		value:          "book",
	}

	nodeJ := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeJ,
	}
	hashNodeJ := nodeJ.hash_node()
	mpt.db[hashNodeJ] = nodeJ

	flagValueNodeH := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeH := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeH,
		branch_value: [17]string{"", "", hashNodeJ, "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	}
	hashNodeH := nodeH.hash_node()
	mpt.db[hashNodeH] = nodeH

	flagValueNodeD := Flag_value{
		encoded_prefix: nil,
		value:          "verb",
	}

	nodeD := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeD,
		branch_value: [17]string{"", "", hashNodeH, "", "", "", hashNodeE, "", "", "", "", "", "", "", "", "", "verb"},
	}
	hashNodeD := nodeD.hash_node()
	mpt.db[hashNodeD] = nodeD

	flagValueNodeB := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15}),
		value:          hashNodeD,
	}

	nodeB := Node{
		node_type:  2, //Extension
		flag_value: flagValueNodeB,
	}
	hashNodeB := nodeB.hash_node()
	mpt.db[hashNodeB] = nodeB

	flagValueNodeC := Flag_value{
		encoded_prefix: compact_encode([]uint8{6, 15, 7, 2, 7, 3, 6, 5, 16}),
		value:          "stallion",
	}

	nodeC := Node{
		node_type:  2, //Leaf
		flag_value: flagValueNodeC,
	}
	hashNodeC := nodeC.hash_node()
	mpt.db[hashNodeC] = nodeC

	flagValueNodeA := Flag_value{
		encoded_prefix: nil,
		value:          "",
	}

	nodeA := Node{
		node_type:    1, //Branch
		flag_value:   flagValueNodeA,
		branch_value: [17]string{"", "", "", "", hashNodeB, "", "", "", hashNodeC, "", "", "", "", "", "", "", ""},
	}
	hashNodeA := nodeA.hash_node()
	mpt.db[hashNodeA] = nodeA

	flagValueRoot := Flag_value{
		encoded_prefix: compact_encode([]uint8{6}),
		value:          hashNodeA,
	}

	nodeRoot := Node{
		node_type:  2, //Extension Root
		flag_value: flagValueRoot,
	}
	hashNodeRoot := nodeRoot.hash_node()
	mpt.db[hashNodeRoot] = nodeRoot
	mpt.root = hashNodeRoot

	////////////////////////////////////////

	//Add another node to MPT

	//mpt.root = Compact_encode(key) //a
	return errors.New("Problem occured while creating Root Node")
}

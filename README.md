# Merkle Patricia Trie

Ethereum uses a Merkle Patricia Tree (Links to an external site.) to store the transaction data in a block. By organizing the transaction data in a Merkle Patricia Tree, any block with fraudulent transactions would not match the tree's root hash. Build a implementation of a Merkle Patricia Trie, following the specifications at the Ethereum wiki (Links to an external site.).

### Project specification
####
For this project, implemented a Merkle Patricia Trie according to this Link (Links to an external site.).

#### five features of the Merkle Patricia Trie:

### 1. Get(key) -> value
Description: The Get function takes a key as argument, traverses down the Merkle Patricia Trie to find the value, and returns it. If the key doesn't exist, it will return an empty string.(for the Go version: if the key is nil, Get returns an empty string.)
##### Arguments: key (string)
##### Return: the value stored for that key (string).
##### Go function definition: func (mpt *MerklePatriciaTrie) Get(key string) string

### 2. Insert(key, value)
Description: The Insert function takes a key and value as arguments. It will traverse  Merkle Patricia Trie, find the right place to insert the value, and do the insertion.(for the Go version: you can assume the key and value will never be nil.)
##### Arguments: key (string), value (string)
##### Return: string
##### Go function definition: func (mpt *MerklePatriciaTrie) Insert(key string, new_value string)

### 3. Delete(key)
Description: The Delete function takes a key as argument, traverses the Merkle Patricia Trie and finds that key. If the key exists, delete the corresponding value and re-balance the trie if necessary, then return an empty string; if the key doesn't exist, return "path_not_found".
##### Arguments: key (string)
##### Return: string
##### Go function definition: func (mpt *MerklePatriciaTrie) Delete(key string) string

### 4. compact_encode()
The compact_encode function takes an array of numbers as input (each number is between 0 and 15 included, representing a single hex digit), and returns an array of numbers according to the compact encoding rules in the github wiki page under "Compact encoding of hex sequence with optional terminator"). Each number in the output is between 0 and 255 included (representing an ASCII-encoded letter, or for the first value it represents the node type as per the wiki page). You may find a Python version in this Link (Links to an external site.), but be mindful that the return type is different!
##### Arguments: hex_array(array of u8)
##### Return: array of u8
##### Example: input=[1, 6, 1], encoded_array=[1, 1, 6, 1], output=[17, 97]

### 5. compact_decode()
Description: This function reverses the compact_encode() function. 
##### Arguments: hex_array(array of u8)
##### Return: array of u8
##### Example: input=[17, 97], output=[1, 6, 1]

#### Other help functions:

##### 1. fn hash_node(node: &Node) -> String
Description: This function takes a node as the input, hash the node and return the hashed string.

Classes specification: 
In this project, there are two important classes.

1. enum Node
This class represent a node of type Branch, Leaf, Extension, or Null.

2. struct MerklePatriciaTrie
This class represent a Merkle Patricia Trie. It has two variables: "db" and "root".
Variable "db" is a HashMap. The key of the HashMap is a Node's hash value. The value of the HashMap is the Node. 
Variable "root" is a String, which is the hash value of the root node.

Other info:
1. Leaf node and Extension node are differentiated by their prefix, not the enum type.

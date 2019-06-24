#include "../../../vntlib/vntlib.h"

KEY string INTERFACE_META_ID = "0x01ffc9a7";
KEY string ADDR_INTERFACE_ID = "0x3b3b57de";
KEY string CONTENT_INTERFACE_ID = "0xd8389dc5";
KEY string NAME_INTERFACE_ID = "0x691f3431";
KEY string ABI_INTERFACE_ID = "0x2203ab56";
KEY string PUBKEY_INTERFACE_ID = "0xc8690233";
KEY string TEXT_INTERFACE_ID = "0x59d1d43c";

EVENT AddrChanged(indexed string node, address a);
EVENT ContentChanged(indexed string node, string hash);
EVENT NameChanged(indexed string node, string name);
EVENT ABIChanged(indexed string node, indexed uint256 contentType);
EVENT PubkeyChanged(indexed string node, string x, string y);
EVENT TextChanged(indexed string node, indexed string indexedKey, string key);

typedef struct {
  string x;
  string y;
} PublicKey;

typedef struct {
  address addr;
  string content;
  string name;
  PublicKey pubkey;
  mapping(string, string) text;
  mapping(uint256, string) abis;
} Record;

KEY address vns;
KEY mapping(string, Record) records;

CALL address owner(CallParams params, string node);

/**
 * Constructor.
 * @param vnsAddr The VNS registrar contract.
 */
constructor PublicResolver(address vnsAddr) { vns = vnsAddr; }

void onlyOwner(string node) {
  CallParams params = {vns, U256(0), 100000};
  address addr = owner(params, node);
  address sender = GetSender();
  if (Equal(addr, sender) != true) {
    Revert("Not Sender");
  }
}

/**
 * Returns true if the resolver implements the interface specified by the
 * provided hash.
 * @param interfaceID The ID of the interface to check for.
 * @return True if the contract implements the requested interface.
 */
UNMUTABLE
bool supportsInterface(string interfaceID) {
  return Equal(interfaceID, ADDR_INTERFACE_ID) ||
         Equal(interfaceID, CONTENT_INTERFACE_ID) ||
         Equal(interfaceID, NAME_INTERFACE_ID) ||
         Equal(interfaceID, ABI_INTERFACE_ID) ||
         Equal(interfaceID, PUBKEY_INTERFACE_ID) ||
         Equal(interfaceID, TEXT_INTERFACE_ID) ||
         Equal(interfaceID, INTERFACE_META_ID);
}

/**
 * Returns the address associated with an VNS node.
 * @param node The VNS node to query.
 * @return The associated address.
 */
MUTABLE
address addr(string node) {
  records.key = node;
  Record record = records.value;
  return record.addr;
}

/**
 * Sets the address associated with an VNS node.
 * May only be called by the owner of that node in the VNS registry.
 * @param node The node to update.
 * @param addr The address to set.
 */
MUTABLE
void setAddr(string node, address addr) {
  onlyOwner(node);
  records.key = node;
  records.value.addr = addr;
  AddrChanged(node, addr);
}

/**
 * Returns the content hash associated with an VNS node.
 * Note that this resource type is not standardized, and will likely change
 * in future to a resource type based on multihash.
 * @param node The VNS node to query.
 * @return The associated content hash.
 */
UNMUTABLE
string content(string node) {
  records.key = node;
  Record record = records.value;
  return record.content;
}

/**
 * Sets the content hash associated with an VNS node.
 * May only be called by the owner of that node in the VNS registry.
 * Note that this resource type is not standardized, and will likely change
 * in future to a resource type based on multihash.
 * @param node The node to update.
 * @param hash The content hash to set
 */
MUTABLE
void setContent(string node, string hash) {
  onlyOwner(node);
  records.key = node;
  records.value.content = hash;
  ContentChanged(node, hash);
}

/**
 * Returns the name associated with an VNS node, for reverse records.
 * Defined in EIP181.
 * @param node The VNS node to query.
 * @return The associated name.
 */
UNMUTABLE
string name(string node) {
  records.key = node;
  Record record = records.value;
  return record.name;
}

/**
 * Sets the name associated with an VNS node, for reverse records.
 * May only be called by the owner of that node in the VNS registry.
 * @param node The node to update.
 * @param name The name to set.
 */
MUTABLE
void setName(string node, string name) {
  onlyOwner(node);
  records.key = node;
  records.value.name = name;
  NameChanged(node, name);
}

/**
 * Returns the ABI associated with an VNS node.
 * Defined in EIP205.
 * @param node The VNS node to query
 * @param contentTypes A bitwise OR of the ABI formats accepted by the caller.
 * @return contentType The content type of the return value
 * @return data The ABI data
 */
UNMUTABLE
string ABIRecord(string node, uint256 contentTypes) {
  records.key = node;
  Record record = records.value;
  for (uint256 contentType = U256(1); U256_Cmp(contentType, contentTypes) != 1;
       contentType = U256_Shl(contentType, 1)) {
    record.abis.key = contentType;
    string recordabis = record.abis.value;
    if (U256_Cmp(U256_And(contentType, contentTypes), U256(0)) != 0 &&
        Equal(recordabis, "") == false) {
      return recordabis;
    }
  }
  return "";
}

UNMUTABLE
uint256 ABIContentType(string node, uint256 contentTypes) {
  records.key = node;
  Record record = records.value;
  for (uint256 contentType = U256(1); U256_Cmp(contentType, contentTypes) != 1;
       contentType = U256_Shl(contentType, 1)) {
    record.abis.key = contentType;
    string recordabis = record.abis.value;
    if (U256_Cmp(U256_And(contentType, contentTypes), U256(0)) != 0 &&
        Equal(recordabis, "") == false) {
      return contentType;
    }
  }
  return U256(0);
}

/**
 * Sets the ABI associated with an VNS node.
 * Nodes may have one ABI of each content type. To remove an ABI, set it to
 * the empty string.
 * @param node The node to update.
 * @param contentType The content type of the ABI
 * @param data The ABI data.
 */
MUTABLE
void setABI(string node, uint256 contentType, string data) {
  onlyOwner(node);
  // Content types must be powers of 2
  if (U256_Cmp(U256_Add(U256_Sub(contentType, 1), contentType), 0) != 0) {
    Revert("");
  }
  records.key = node;
  records.value.abis.key = contentType;
  records.value.abis.value = data;
  ABIChanged(node, contentType);
}

/**
 * Returns the SECP256k1 public key associated with an VNS node.
 * Defined in EIP 619.
 * @param node The VNS node to query
 * @return x, y the X and Y coordinates of the curve point for the public key.
 */
UNMUTABLE
string pubkeyX(string node) {
  records.key = node;
  return records.value.pubkey.x;
}

UNMUTABLE
string pubkeyY(string node) {
  records.key = node;
  return records.value.pubkey.y;
}

/**
 * Sets the SECP256k1 public key associated with an VNS node.
 * @param node The VNS node to query
 * @param x the X coordinate of the curve point for the public key.
 * @param y the Y coordinate of the curve point for the public key.
 */
MUTABLE
void setPubkey(string node, string x, string y) {
  onlyOwner(node);
  records.key = node;
  records.value.pubkey.x = x;
  records.value.pubkey.y = y;
  PubkeyChanged(node, x, y);
}

/**
 * Returns the text data associated with an VNS node and key.
 * @param node The VNS node to query.
 * @param key The text data key to query.
 * @return The associated text data.
 */
UNMUTABLE
string text(string node, string key) {
  records.key = node;
  records.value.text.key = key;
  return records.value.text.value;
}

/**
 * Sets the text data associated with an VNS node and key.
 * May only be called by the owner of that node in the VNS registry.
 * @param node The node to update.
 * @param key The key to set.
 * @param value The text data value to set.
 */
MUTABLE
void setText(string node, string key, string value) {
  onlyOwner(node);
  records.key = node;
  records.value.text.key = key;
  records.value.text.value = value;
  TextChanged(node, key, key);
}

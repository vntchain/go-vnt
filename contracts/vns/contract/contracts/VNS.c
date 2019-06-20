#include "../../../vntlib/vntlib.h"

typedef struct {
  address owner;
  address resolver;
  uint64 ttl;
} Record;

KEY mapping(string, Record) records;

// Logged when the owner of a node assigns a new owner to a subnode.
EVENT NewOwner(indexed string node, indexed string label, address owner);

// Logged when the owner of a node transfers ownership to a new account.
EVENT Transfer(indexed string node, address owner);

// Logged when the resolver for a node changes.
EVENT NewResolver(indexed string node, address resolver);

// Logged when the TTL of a node changes
EVENT NewTTL(indexed string node, uint64 ttl);

/**
 * Constructs a new VNS registrar.
 */
constructor VNS() {
  records.key = "";
  records.value.owner = GetSender();
}

/**
 * Returns the address that owns the specified node.
 */
// function owner(bytes32 node) constant returns (address) {
//     return records[node].owner;
// }

void onlyOwner(string node) {
  records.key = node;
  address owner = records.value.owner;
  address sender = GetSender();
  if (Equal(owner, sender) == false) {
    Revert("need owner");
  }
}

UNMUTABLE
address owner(string node) {
  records.key = node;
  return records.value.owner;
}

/**
 * Returns the address of the resolver for the specified node.
 */
// function resolver(bytes32 node) constant returns (address) {
//     return records[node].resolver;
// }

UNMUTABLE
address resolver(string node) {
  records.key = node;
  return records.value.resolver;
}

/**
 * Returns the TTL of a node, and any records associated with it.
 */
// function ttl(bytes32 node) constant returns (uint64) {
//     return records[node].ttl;
// }
UNMUTABLE
uint64 ttl(string node) {
  records.key = node;
  return records.value.ttl;
}

/**
 * Transfers ownership of a node to a new address. May only be called by the
 * current owner of the node.
 * @param node The node to transfer ownership of.
 * @param owner The address of the new owner.
 */
MUTABLE
void setOwner(string node, address owner) {
  onlyOwner(node);
  Transfer(node, owner);
  records.key = node;
  records.value.owner = owner;
}

/**
 * Transfers ownership of a subnode SHA3(node, label) to a new address. May only
 * be called by the owner of the parent node.
 * @param node The parent node.
 * @param label The hash of the label specifying the subnode.
 * @param owner The address of the new owner.
 */
MUTABLE
void setSubnodeOwner(string node, string label, address owner) {
  onlyOwner(node);
  string subnode = SHA3(Concat(node, label));
  NewOwner(node, label, owner);
  records.key = subnode;
  records.value.owner = owner;
}

/**
 * Sets the resolver address for the specified node.
 * @param node The node to update.
 * @param resolver The address of the resolver.
 */
MUTABLE
void setResolver(string node, address resolver) {
  onlyOwner(node);
  NewResolver(node, resolver);
  records.key = node;
  records.value.resolver = resolver;
}

/**
 * Sets the TTL for the specified node.
 * @param node The node to update.
 * @param ttl The TTL in seconds.
 */
MUTABLE
void setTTL(string node, uint64 ttl) {
  onlyOwner(node);
  records.key = node;
  records.value.ttl = ttl;
}
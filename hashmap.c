#include <stdint.h>
#include <malloc.h>
#include <stdio.h>

typedef union {
    double value;
    uint8_t bytes[8];
} double_u;

unsigned int rotateLeft(unsigned int n, int amount) {
    return (n << amount) | (n >> ((sizeof(n) * 8)-amount));
}

unsigned int rotateRight(unsigned int n, int amount) {
    return (n >> amount) | (n << ((sizeof(n) * 8)-amount));
}

#define HASH_SALT 0x79812793

unsigned int hashFunction(double value) {
    double_u number = {.value = value};
    uint8_t *data = number.bytes;
    unsigned int result = 0;
    for(unsigned int i = 0; i < 8; i++) {
        unsigned int modify;
        if(data[i] & 1) modify = rotateRight(rotateLeft(data[i]+i, 3), i*8);
        else modify = rotateLeft(rotateRight(data[i]+i, 3), i*5);
        result ^= modify ^ (modify >> i*11) ^ (modify << i*3);
    }
    result += HASH_SALT;
    return result ^ rotateRight(result, 16);
}

#define HASHMAP_BUCKET_MASK 0xFFFF
#define HASHMAP_BUCKET_INITSIZE 8

typedef struct {
    double key;
    double value;
} bucketItem_t;

typedef struct {
    bucketItem_t *items;
    uint32_t length;
    uint32_t itemCount;
} bucket_t;

typedef struct {
    bucket_t *buckets;
    unsigned int itemCount;
} hashmap_t;

hashmap_t *hashmapCreate() {
    hashmap_t *hashmap = calloc(1, sizeof(hashmap_t));
    hashmap->buckets = calloc(HASHMAP_BUCKET_MASK+1, sizeof(bucket_t));
    return hashmap;
}

double hashmapRead(hashmap_t *hashmap, double key) {
    unsigned int hash = hashFunction(key) & HASHMAP_BUCKET_MASK;

    //Is bucket unitinialized?
    bucket_t *bucket = &(hashmap->buckets[hash]);
    if(bucket->items == NULL) {
        return key;
    }

    //Loop through keys
    for(unsigned int i = 0; i < bucket->itemCount; i++) {
        if(bucket->items[i].key == key) {
            return bucket->items[i].value;
        }
    }

    //No key matched, return key itself
    return key;
}

void hashmapDelete(hashmap_t *hashmap, double key) {
    unsigned int hash = hashFunction(key) & HASHMAP_BUCKET_MASK;
    bucket_t *bucket = &(hashmap->buckets[hash]);
    
    //Is bucket no initialized?
    if(bucket->items == NULL) {
        return;
    }

    //Find key in bucket
    for(unsigned int i = 0; i < bucket->itemCount; i++) {
        if(bucket->items[i].key == key) {
            bucket->itemCount--;
            if(bucket->itemCount == 0) {
                free(bucket->items);
            }
            else {
                bucketItem_t *item = &(bucket->items[bucket->itemCount]);
                bucket->items[i].key = item->key;
                bucket->items[i].value = item->value;
            }
        }
    }

    //No key found
    return;
}

void hashmapWrite(hashmap_t *hashmap, double key, double value) {
    if(key == value) {
        hashmapDelete(hashmap, key);
        return;
    }
    unsigned int hash = hashFunction(key) & HASHMAP_BUCKET_MASK;
    
    //Is bucket unitinialized?
    bucket_t *bucket = &(hashmap->buckets[hash]);
    if(bucket->items == NULL) {
        bucket->items = calloc(HASHMAP_BUCKET_INITSIZE, sizeof(bucketItem_t));
        bucket->length = HASHMAP_BUCKET_INITSIZE;
    }

    //Find matching key to overwrite
    unsigned int i = 0;
    for(; i < bucket->itemCount; i++) {
        if(bucket->items[i].key == key) {
            bucket->items[i].value = value;
            return;
        }
    }

    //Resize if bucket is full
    if(bucket->itemCount == bucket->length) {
        bucket->items = realloc(bucket->items, bucket->length*2);
        bucket->length = bucket->length * 2;
    }

    //Place key at the end of bucket
    bucket->items[i].key = key;
    bucket->items[i].value = value;
    bucket->itemCount++;
}

int main() {
    return 0;
}
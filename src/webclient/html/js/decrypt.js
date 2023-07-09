// Decryption
function decryptFile(encodedStruct, encryptedData) {
    // Decode the Base64 encoded struct
    const decodedStruct = JSON.parse(atob(encodedStruct));
  
    // Import the key from the struct
    const importedKey = crypto.subtle.importKey(
      'jwk',
      decodedStruct.key,
      { name: 'AES-GCM', length: 256 },
      true,
      ['decrypt']
    );
  
    // Convert the IV to Uint8Array
    const iv = new Uint8Array(decodedStruct.iv);
  
    // Decrypt the file data
    const decryptedData = crypto.subtle.decrypt(
      { name: 'AES-GCM', iv: iv },
      importedKey,
      encryptedData
    );
  
    return decryptedData;
  }
  
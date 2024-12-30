<template>
    <div>
        <h2>Wallet Upload</h2>
        <form @submit.prevent="processFiles">
            <div>
                <label for="cert">Upload Certificate:</label>
                <input type="file" id="cert" @change="onFileChange('cert', $event)" required />
            </div>
            <div>
                <label for="pkey">Upload Private Key:</label>
                <input type="file" id="pkey" @change="onFileChange('pkey', $event)" required />
            </div>
            <button type="submit">Process Files</button>
        </form>
        <p v-if="responseMessage">{{ responseMessage }}</p>
    </div>
</template>

<script>
import api from "../api/api";

export default {
    data() {
        return {
            files: {
                cert: null,
                pkey: null,
            },
            privateKey: null,
            certificate: null,
            responseMessage: null,
        };
    },
    methods: {
        onFileChange(fileKey, event) {
            const file = event.target.files[0];
            this.files[fileKey] = file;

            const regex =
                fileKey === "pkey"
                    ? /-----BEGIN PRIVATE KEY-----(.*?)-----END PRIVATE KEY-----/s
                    : /-----BEGIN CERTIFICATE-----(.*?)-----END CERTIFICATE-----/s;

            this.extractKey(file, fileKey === "pkey" ? "privateKey" : "certificate", regex);
        },
        extractKey(file, keyType, regex) {
            const reader = new FileReader();
            reader.onload = (event) => {
                const fileContent = event.target.result;
                const match = fileContent.match(regex);
                console.log(`Extracting ${keyType}:`, match);
                if (match && match[0]) {
                    // Extract the full PEM block and remove BEGIN/END lines
                    const cleanedKey = match[0]
                        .replace(/-----BEGIN [^-]+-----/g, "")
                        .replace(/-----END [^-]+-----/g, "")
                        .replace(/\n/g, "");

                    this[keyType] = cleanedKey.trim();
                    console.log(`${keyType} extracted successfully.`);
                } else {
                    this.responseMessage = `Invalid ${keyType === "privateKey" ? "private key" : "certificate"} format.`;
                }
            };
            reader.onerror = (error) => {
                console.error(`Error reading ${keyType} file:`, error);
                this.responseMessage = `Failed to read the ${keyType} file.`;
            };
            reader.readAsText(file);
        },

        async processFiles() {
            if (!this.privateKey || !this.certificate) {
                this.responseMessage = "Please upload both valid files.";
                return;
            }

            try {
                const requestBody = {
                    privateKey: this.privateKey,
                    certificate: this.certificate,
                };

                const response = await api.post("/wallet_sign_in", requestBody, {
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                this.responseMessage = `Files processed successfully: ${response.data.message}`;
            } catch (error) {
                console.error("Error processing files:", error);
                this.responseMessage = "Failed to process files.";
            }
        },
    },
};
</script>
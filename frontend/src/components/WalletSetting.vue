<template>
    <div>
        <h2>Wallet Sign</h2>
        <form @submit.prevent="processFiles">
            <div>
                <label for="cert">Upload Certificate:</label>
                <input type="file" id="cert" @change="onFileChange('cert', $event)" required />
            </div>
            <div>
                <label for="pkey">Upload Private Key:</label>
                <input type="file" id="pkey" @change="onFileChange('pkey', $event)" required />
            </div>
            <div>
                <label for="msp">Upload MSP File:</label>
                <input type="file" id="msp" @change="onFileChange('msp', $event)" required />
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
                msp: null,
            },
            privateKey: null,
            certificate: null,
            mspContent: null,
            responseMessage: null,
        };
    },
    methods: {
        onFileChange(fileKey, event) {
            const file = event.target.files[0];
            this.files[fileKey] = file;

            const regexMap = {
                pkey: /-----BEGIN PRIVATE KEY-----(.*?)-----END PRIVATE KEY-----/s,
                cert: /-----BEGIN CERTIFICATE-----(.*?)-----END CERTIFICATE-----/s,
                msp: /.+/,
            };

            this.extractFileContent(file, fileKey, regexMap[fileKey]);
        },
        extractFileContent(file, keyType, regex) {
            const reader = new FileReader();
            reader.onload = (event) => {
                const fileContent = event.target.result;
                const match = fileContent.match(regex);

                if (match) {
                    this[keyType === "pkey" ? "privateKey" : keyType === "cert" ? "certificate" : "mspContent"] = fileContent.trim();
                    console.log(`${keyType} content extracted successfully.`);
                } else {
                    this.responseMessage = `Invalid ${keyType.toUpperCase()} format or content.`;
                }
            };
            reader.onerror = (error) => {
                console.error(`Error reading ${keyType} file:`, error);
                this.responseMessage = `Failed to read the ${keyType.toUpperCase()} file.`;
            };
            reader.readAsText(file);
        },

        async processFiles() {
            if (!this.privateKey || !this.certificate || !this.mspContent) {
                this.responseMessage = "Please upload all required files.";
                return;
            }

            try {
                const requestBody = {
                    privateKey: this.privateKey,
                    certificate: this.certificate,
                    mspContent: this.mspContent,
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
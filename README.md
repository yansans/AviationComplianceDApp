# AviationComplianceDApp

## Deskripsi Aplikasi

Aplikasi AviationComplianceDApp dibuat sebagai aplikasi verifikasi dan konfirmasi terhadap compliance suatu perusahaan penerbangan dengan aturan serta sebagai solusi dari tantangan kebergantungan terhadap perusahaan, risiko manipulasi, dan transparansi proses. Use case dari aplikasi ini antara lain membuat compliance asset, membaca compliance asset, meng-update compliance asset, serta mengakses riwayat compliance asset. Tech Stack yang digunakan adalah Hyperledger Fabric sebagai platform blockchain, Fabric Gateway sebagai backend, Vue sebagai frontend, serta Python sebagai Oracle.

## Requirement

Rekomendasi menggunakan Linux atau WSL

1. Docker
2. Go
3. Fabric Binary dan Docker (Panduan lebih lanjut dalam readme di folder fabric)

## Cara Menjalankan Private Chain

1. Clone repository `git clone https://github.com/yansans/AviationComplianceDApp.git`
2. Masuk ke folder fabric `cd fabric` kemudian masuk ke folder test-network `cd test-network`
3. Jalankan script network `./network down` kemudian `./network.sh up createChannel -c channel1`

## Cara Deployment Smart Contract

1. Masuk ke folder fabric `cd fabric` kemudian masuk ke folder test-network `cd test-network`
2. Jalankan script smart contract `./network.sh deployCC -ccn basic -ccp ../../backend/chaincode/ -ccl go -c channel1`

## Cara Deployment dan Integrasi Oracle

0. Pastikan .env sudah terisi dengan benar
1. Masuk ke folder backend `cd backend` kemudian masuk ke folder api `cd api`
2. Jalankan backend yang secara langsung akan menjalankan oracle `go run .`

## Cara menjalankan frontend

1. Masuk ke folder frontend dengan `cd frontend`
2. Install dependency dengan `npm install`
3. Jalankan frontend dengan `npm run dev`

## Link Video Demonstrasi

[Video Demo](https://drive.google.com/file/d/1SstqAB5AFBSbh2hc_0__7JzNcfeyA5CW/view?usp=sharing)

## 🔧 Installation

### 🛠 Prerequisites

Make sure you have the following installed:

- **Go ≥ 1.21**
- **C++17** compiler (`g++`)
- **NTL** and **GMP** development libraries

Install on Ubuntu/Debian:

```
sudo apt update
sudo apt install build-essential libntl-dev libgmp-dev
```

---

### ⚙️ Step 1: Compile Native C++ Libraries

Run the setup script to compile the required C++ components:

```
./run.sh
```

This generates the following static libraries:

- `tdm/libNTT.a`
- `ecc/libReedSolomon.a`
- `mvp/libMVP.a`

**Note:** You must run this command from the **root** of the repository.

---

### 🌐 Step 2: Export Paths for Go to Use Native Libraries

Before running tests, set the appropriate environment variables **from the root directory** of this project:

```
export CGO_CXXFLAGS="-I$(pwd)/tdm -I$(pwd)/ecc -I/usr/include"
export CGO_LDFLAGS="-L$(pwd)/tdm -L$(pwd)/ecc -lNTT -lReedSolomon -lntl -lgmp"
```

---

### 🚀 Step 3: Run Go Benchmarks

From the project root (or inside the `mvp/` folder if that's where your `go.mod` is), run:

```
go test -bench=. ./...
```

This will compile the Go code, link against the native C++ libraries, and execute the performance benchmarks.

---

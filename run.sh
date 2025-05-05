g++ -std=c++17 -O3 -c tdm/NTT.cpp -o tdm/NTT.o -I/opt/homebrew/include
ar rcs tdm/libNTT.a tdm/NTT.o

g++ -std=c++17 -O3 -c mvp/mvp.cpp -o mvp/mvp.o -I/opt/homebrew/include
ar rcs mvp/libMVP.a mvp/mvp.o

g++ -std=c++17 -O3 -c ecc/ReedSolomon.cpp -o ecc/ReedSolomon.o -I/opt/homebrew/include
ar rcs ecc/libReedSolomon.a ecc/ReedSolomon.o

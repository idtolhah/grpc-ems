#!/bin/bash
cd bff
skaffold dev
cd ..

cd comment-cmd
skaffold dev
cd ..

cd master-query
skaffold dev
cd ..

cd packing-cmd
skaffold dev
cd ..

cd packing-query
skaffold dev
cd ..

cd user-query
skaffold dev
cd ..

# wait $P1 $P2 $P3 $P4 $P5 $P6
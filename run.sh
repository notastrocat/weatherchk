#!/bin/bash

export API_KEY=MVJ2XRJEGFPXG3FAP6SCPTCGZ
export COLORTERM=truecolor

# command to start the app...
executable=$(ls | grep weatherchk)

build() {
  build_err=$(go build -o weatherchk 2>&1 | tee build.log; grep main.go build.log) 
  rm build.log
}

if [[ $executable == "" ]]; then
  echo "executable present: NO"
  echo "building project..."
  # build_err=$(go build -o weatherchk 2>&1 | tee build.log; grep main.go build.log)
  # rm build.log
  build
else
  echo "executable present: YES ($executable)"

  if [[ $1 == "rebuild" ]]; then
    echo "found rebuild flag. Rebuilding the app."
    rm $executable
    build
  fi
fi

if [[ $build_err == "" ]]; then
  echo "no build errors encountered; running the application..."
  ./weatherchk
else
  echo "build errors encountered!!!..."
  echo $build_err
fi


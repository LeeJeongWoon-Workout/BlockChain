package main

import (
    "fmt"
    "math"
)

func is_prime(n int) bool {
    result:=true

    for i:=2;i<=int(math.Sqrt(float64(n)));i++ {
        if n%i==0 {
            result=false
        }
    }
    return result
}

func primes(begin,end int) []int {
    v:=[]int{}

    for i:=begin;i<=end;i++ {
        if is_prime(i) {
            v=append(v,i)
        }
    }
    return v
}

func main() {
    v:=primes(1000,10000)
    fmt.Println(v)
}
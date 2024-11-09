# Benchmark comparison

This sample will compare Ore (current commit of Nov 2024) to [samber/do/v2 v2.0.0-beta.7](https://github.com/samber/do).
We registered the below dependency graphs to both Ore and SamberDo, then ask them to create the concrete `A`.

We will only benchmark the creation, not the registration. Because registration usually happens only once on application startup =>
 not very interesting to benchmark.

## Data Model

- This data model has only 2 singletons `F` and `Gb` => they will be created only once
- Other concretes are `Transient` => they will be created each time the container create a new `A` concrete.
- We don't test the "Scoped" lifetime in this excercise because SamberDo doesn't have equivalent support for it. [The "Scoped" functionality of SamberDo](https://do.samber.dev/docs/container/scope) means "Sub Module" rather than a lifetime.

```mermaid
flowchart TD
A["A<br><sup></sup>"]
B["B<br><sup></sup>"]
C["C<br><sup></sup>"]
D["D<br><sup><br></sup>"]
E["E<br><sup><br></sup>"]
F["F<br><sup>Singleton</sup>"]
G(["G<br><sup>(interface)</sup>"])
Gb("Gb<br><sup>Singleton</sup>")
Ga("Ga<br><sup></sup>")
DGa("DGa<br><sup>(decorator)</sup>")
H(["H<br><sup>(interface)<br></sup>"])
Hr["Hr<br><sup>(real)</sup>"]
Hm["Hm<br><sup>(mock)</sup>"]

A-->B
A-->C
B-->D
B-->E
D-->H
D-->F
Hr -.implement..-> H
Hm -.implement..-> H
E-->DGa
E-->Gb
E-->Gc
DGa-->|decorate| Ga
Ga -.implement..-> G
Gb -.implement..-> G
Gc -.implement..-> G
DGa -.implement..-> G
```

## Run the benchmark by yourself

```sh
 go test -benchmem -bench .
 ```

## Sample results

On my machine, Ore always perform faster and use less memory than Samber/Do:

```text
Benchmark_Ore-12                  415822              2565 ns/op            2089 B/op         57 allocs/op
Benchmark_SamberDo-12             221941              4954 ns/op            2184 B/op         70 allocs/op
```

And with `ore.DisableValidation = true`

```text
Benchmark_Ore-12                  785088              1668 ns/op            1080 B/op         30 allocs/op
Benchmark_SamberDo-12             227851              4940 ns/op            2184 B/op         70 allocs/op
```

As any benchmarks, please take these number "relatively" as a general idea:

- These numbers are probably outdated at the moment you are reading them
- You might got a very different numbers when running them on your machine or on production machine.

# k8autotune - Automatic Tuner for Kubernetes

This service evaluates POD parameters to determine a better configuration to improve the resource use reducing the total number of nodes a cluster will need.

The current *resources* spec is taken from the POD controller and K8AutoTune will monitor the POD activity so both values will be checked adjusting their values if there is a big gap between both values.

Different strategies will be allowed to decide the best value for each use profile.

# Roadmap

## Version 0.1

* [ ] Handle Deployments
* [ ] Measure CPU activity
* [ ] Automatically adjust CPU ussage
* [ ] Use Prometheus as metrics source
* [ ] Update limit if needed

## Pending
* [ ] Allow to configure namespaces to include
* [ ] Persist metrics ussage locally in file system
* [ ] Automatically adjust Memory ussage
* [ ] Allow to annotate objects to configure behaviour
* [ ] Configure by custom resource
* [ ] Emit allert of changes
* [ ] Emit allert of exceptional conditions (like possible memory leaks)
* [ ] Take into account warm-up times for heavy PODs (like those that use a JVM)
* [ ] Let one controller to stabilize before tuning next
* [ ] 

# References

Here many sources with good ideas many of them used in this project

* [Object-Oriented Software Engineering: A Use Case Driven Approach](https://www.amazon.com/Object-Oriented-Software-Engineering-Driven-Approach/dp/0201403471)
* [Agile Software Development, Principles, Patterns, and Practices](https://www.amazon.com/dp/0135974445/ref=wl_it_dp_o_pC_nS_ttl?_encoding=UTF8&colid=CG11VVP0H8Y8&coliid=I1P9T8D1QRUFMM)
* [DDD, Hexagonal, Onion, Clean, CQRS, … How I put it all together](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/)
* 
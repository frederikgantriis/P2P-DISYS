# Description:

You have to implement distributed mutual exclusion between nodes in your distributed system. 

Your system has to consist of a set of peer nodes, and you are not allowed to base your implementation on a central server solution.

You can decide to base your implementation on one of the algorithms, that were discussed in lecture 7.

# System Requirements:

> - R1: Implement a system with a set of peer nodes, and a Critical Section, that represents a sensitive system operation. Any node can at any time decide it wants access to the Critical Section. Critical section in this exercise is emulated, for example by a print statement, or writing to a shared file on the network.
> - R2: Safety: Only one node at the same time is allowed to enter the Critical Section 
> - R2: Liveliness: Every node that requests access to the Critical Section, will get access to the Critical Section (at some point in time)

# Technical Requirements:

1. Use Golang to implement the service's nodes
1. Provide a README.md, that explains how to start your system
1. Use gRPC for message passing between nodes
1. Your nodes need to find each other.  For service discovery, you can choose one of the following options
1. Supply a file with IP addresses/ports of other nodes
1. Enter IP address/ports through the command line use a package for service discovery, like the Serf package 
1. Demonstrate that the system can be started with at least 3 nodes
1. Demonstrate using your system's logs,  a sequence of messages in the system, that leads to a node getting access to the Critical Section. You should provide a discussion of your algorithm, using examples from your logs.

# Hand-in requirements:

>1. Hand in a single report in a pdf file
>1. Provide a link to a Git repo with your source code in the report
>1. Include system logs, that document the requirements are met, in the appendix of your report

# Grading notes

Partial implementations may be accepted, if the students can reason what they should have done in the report.
In order to pass, the students have to attempt to answer all questions.

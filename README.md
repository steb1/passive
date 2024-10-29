## Introduction

Information gathering is a fundamental phase in penetration testing and can be one of the most time-consuming tasks. This project aims to create a passive reconnaissance tool that will help you develop open-source investigative skills.

# Objective

The primary goal of this project is to become more familiar with open-source investigative methods and to develop a tool that collects specific information based on user input.

# Features

This tool allows the user to input data in the form of:

- Full Name (formatted as "Last Name" and "First Name"),
- IP Address,
- Username (e.g., @login format).

Based on the input type, the tool will perform various tasks to gather relevant information.

# Functional Guidelines

## Full Name Input
If the input is a full name, the tool should identify both last and first names and then search available directories for any associated telephone numbers and addresses.

## Address Input
If an IP address is provided, the tool should retrieve and display at least the associated city and Internet Service Provider (ISP) name.

## Username Input
If a username is entered, the tool will search across at least five popular social networks to check if the username exists.

## Output
Results will be saved in a file named result.txt. If result.txt already exists, results will be stored in result2.txt.

# Example Usage
## build project
`go build`

## Example 1: Full Name Input
`./passive --fn "Jean Dupont"`

## IP Address Lookup
`./passive --ip 127.0.0.1`

## Username Search
`./passive --u "@user01"`

# 📑 Technical Requirements

Language Choice: Golang.

🤝 Contributing
Contributions are welcome! Feel free to open issues or submit pull requests to improve the tool.











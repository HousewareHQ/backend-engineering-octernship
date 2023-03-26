# Organization's API (Authentication+Authorization)

## **Please Refer README.md files in <br> '/backend' - Setting up backend and frontend. <br> '/backend/api'-Getting started with testing API. <br> '/backend/web/app'-Getting started with web client.<br> Directories to get started with setting up and testing the assignment project.**
© 2023 github.com/thepranays. All rights reserved.

---


---







[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-c66648af7eb3fe8bc4f294546bfd86ef473780cde1dea487d3c4ff354943c9ae.svg)](https://classroom.github.com/online_ide?assignment_repo_id=10472174&assignment_repo_type=AssignmentRepo)
## Houseware
### Company information 

Houseware's vision is to empower the next generation of knowledge workers by putting the data warehouse in their hands, in the language they speak. Houseware is purpose-built for the Data Cloud’s untouched creators, empowering internal apps across organizations. 

### Why participate in an Octernship with Houseware

Houseware is changing the way the data warehouse is leveraged, and we want you to help build Houseware! Our team came together to answer the singular question, "how can we flip the value of the data warehouse to the ones who really need it, to the ones who drive decisions". 

In this role, you'll have the opportunity to work as a Backend engineer with the Houseware team on multiple customer-facing projects, the role being intensive in technical architecture and backend engineering. The ability to have a constant pulse on the engineering team’s shipping velocity, while accounting for stability and technical debt looking forward is crucial.

### Octernship role description

We're looking for backend developers to join the Houseware team. 

We are hell-bent on building a forward-looking product, something that constantly pushes us to think by first principles and question assumptions, building a team that is agile in adapting and ever curious. While fast-paced execution is one of the prerequisites in this role, equally important is the ability to pause and take stock of where product/engineering is heading from a long-term perspective. Your initiative is another thing that we would expect to shine through here, as you continuously navigate through ambiguous waters while working with vigor on open-ended questions - all to solve problems for and empathize with the end users.

You are expected to own the backend and infrastructure stack end-to-end, understand the business use cases, map it to the best-in-class engineering systems while maintaining a great developer experience. This role involves a high level of attention to detail, debugging and testing skills, as well as long-term thinking with respect to the scalability of our platform. 


| Octernship info  | Timelines and Stipend |
| ------------- | ------------- |
| Assignment Deadline  | 26 March 2023  |
| Octernship Duration  | 3-6 Months  |
| Monthly Stipend  | $600 USD  |

### Recommended qualifications

You’d be a great fit if:

- You’re proficient in Golang and Python, having prior experience building backend systems and hands-on experience with AWS/GCP.
- You’re familiar with the modern data stack and have a good understanding of Infrastructure-as-code tooling like Terraform.
- Plus Points if you’re a contributor to open-source, we’d love to see your work!

### Eligibility

To participate, you must be:

* A [verified student](https://education.github.com/discount_requests/pack_application) on Global Campus

* 18 years or older

* Active contributor on GitHub (monthly)

# Assignment

## Implement an Authorization+Authentication service in Golang

### Task instructions

The assignment is to create a backend API service in Golang that handles authorization and authentication for a web app. The details of the web app are as follows:
- A simple web app where users in an organization can signin and list all other users in their organization
- Logging in is performed by supplying a `username, password` combination
- Note that all passwords should be hashed when stored in a database for security purposes
- For simplicity, assume that the existing users have already been registered and we are not concerned about a user registration flow here.
- The user should be logged in with a JWT token, with a one hour expiry.
- The user should be able to receive a new access token using a 'Refresh token' with a validity of 24 hours.
- The user should be able to logout as well.
- There are admin privileges assigned to a few users, which gives them the ability to add new user accounts or delete existing user accounts from their organization.
- All non-admin users should be able to see other user accounts but shouldn't be able to add/delete any user accounts.
- Note that any user shouldn't be able to view/add/delete user accounts into any other organization.

The API should follow REST API conventions, feel free to design the API structure as you may. The API should cover the following functionalities:
- User Login
- User Logout
- Admin User adds a new User account(by providing the username & password)
- Admin User deletes an existing User account from their organization
- List all Users in their organization

Note: Do add unit tests(for success & failure) for each API endpoint.

Provided in this Github template is a Golang Standard repository, you'd have to design an ideal architecture/stack for this problem
- Golang framework for this API
- Which Database shall be used to store the user details?
- The ORM that shall be used for interfacing with the Database
- DB design

Do document the design decisions and the rationale behind the same in a README file.

1. Please push your final code changes to your main branch
2. Please add your instructions/assumptions/any other remarks in the beginning of the Readme file and the reviewers will take a look
3. The PR created called Feedback will be used for sharing any feedback/asking questions by the reviewers, please make sure you do not close the Feedback PR.
4. The assignment will be automatically submitted on the "Assignment Deadline" date -- you don't need to do anything apart from what is mentioned above.
5. Using GitHub Issues to ask any relevant questions regarding the project


### Task Expectations

- Instructions in the Readme to setup the API & the relevant database
- Postman/Swagger/OpenAPI spec so that the APIs can be tested
- The task will be evaluated on the: fulfillment of the requirements and correctness of the API responses, in addition to the simplicity & architecture of the solution

### Task submission

Students are expected to use the [GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow) when working on their project. 

1. Making changes on the auto generated `feedback` branch to complete the task
2. Using the auto generated **Feedback Pull Request** for review and submission
3. Using GitHub Discussions to ask any relevant questions regarding the project

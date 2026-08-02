This is a general coding guideline's template which should be used to generate a `commands/coding-guideline.md`.

## Package design

The package design of an application must follow these things:

1. the minimum main package, which usually initializes a global object and start a long running process or main process for short-lived process
2. packages for each boundary context. For example, DB client to handle a connection as well as sending SQL queries, HTTP server's handlers to receive HTTP requests, gRPC servers' handler to handle gRPC's request, or 3rd party vendors' sdk wrapper to access 3rd party's vendors' external requests
3. packages for domain layers

The 2nd makes the code reuse possible across specific protocols and also allow to hold global objects easily.
Besides, some of these components can be reused across multiple applications.


## Simplicity and code readability

Simplicity is the most important thing for source code because simple code is the easiest to read the code and it's intent.
There are many levels of techniques we should follow for simplicity and code readability.
Reusing codes are very frequent operations and most core tasks to keep codebase simple and readable.
There are a few things to make code reusable
1. It's important to implement code with high cohesion. Each function should do one thing and when the code can be reused, but many codes do not require that new feature, then consider to define extract new function and reuse the same function without code duplication, but add new feature to the new function, so it can be applied only where it's necessary.
2. Avoid premature optimization. It's easy to implement abstracted code, but abstract codes only when there are 3 or more concrete use cases for simplicity.

Not only reusing codes, but also the code itself should be usable without comments. Every line of code should not have comment on each line. Instead, add the reasons of the code, or only high-level contexts of architectures, designs, or purposes of comments in a block of code or function.
The names of variables or functions should be concise but easy enough to understand for people to understand what they are, and if they are not concise, think about refactoring. Exceptions are variables with local scopes used within around 5 lines of code. This short lived variables can be defined in a very short name for readability like `i` or `p`.
Related to namings, define constants from magic numbers or strings so an intent is clear by the constant name.
Besides, avoid very general names for any names in variable, functions, classes, or packages. For instance, shared, common, utils, or info. These general names can be used in many places, but it blurs responsibilities of each and make the code harder to read, define in a wrong place and fail to reuse.


Of course, it's important not to add unnecessary changes or delete unnecessary old code.
For example, avoid adding codes of validations everywhere for defensive programming. This is opinionated, but it's important to add validation on boundary contexts layers, and remove unnecessary defensive code other than boundary contexts, to ensure both error cases, security, and also simplicity.
For example, don't add validation of all arguments in all functions. Instead, add validation for requests of a HTTP server to make sure they are expected values. Or if an application makes a call to a 3rd party vendor, then make sure the responses of the 3rd party vendor is expected values.
The data in the RDB is a difficult decision, but as long as an application manages that DBs, validation of records are not necessary.
If the code is in a library, make sure to add a validation in public functions so each client code can make sure it works as expected, but don't validate in private functions.
Second, in many cases, it's important to recognize in which case breaking changes are not acceptable. It's important to remove unnecessary and unused code for private functions or alpha version APIs that nobody use for simplicity.
Third, deleting unnecessary codes simply makes code simple and easy to read. This should be performed especially before refactoring such that reusable codes can be found effectively and efficiently. So keep deleting unused functions, branches, imports or variables.


There are a lot of other small techniques that can keep code simple and readable.
1. Separate functions over boolean arguments: when a boolean argument is passed to a function and inside the function, sometimes, it is used to run different logic almost entirely. In this case, this method should be split into multiple functions.
2. it's important to avoid too long arguments like or return types of a function, like more than 5 arguments or more than 3 return types of a function. Instead, define custom types or objects and use them for arguments of return types. This avoid to update callers when there is new change.
3. Early returns over nesting: Use guard clauses by return in a if clause at the top of a block or a function. This is to keep the indent of code as shallow as possible and simple.
4. Reduce exposed variables, functions, or any other objects as much as possible. This allows applications or libraries to be used without intentions and reduce backward compatibilities as much as possible.


### Refactoring process

Refactoring is very important for increasing reusability of source codes for same functionality across the codebase and it makes codebase easy to read and simple. There are a few ways to achieve refactoring. It's important to follow both processes

1. Perform refactoring before adding new changes. Before starting implementing new changes, understand existing codes, figuring out if there are codes that can be reused, refactor some codes that can be reused for new changes, and finally add new changes.
2. Perform refactoring after QA for new changes pass: new changes added must work as expected, but changes originally planned are not perfect in most cases. So after all QA is passed, it's important to make sure refactoring should be conducted for code reusability.


## Testability

The most important thing for your code is easy to do testings.
The tests must be deterministic, must not be flakey, and achieve at least 95% branch coverage. In order to achieve high testability, there are many things we have to ensure for non-testing codes, i.e. production code
First, we have to design the code to inject dependencies easily by constructor or method injections. For example, it's very easy to write non deterministic code when time is related. sleep function must be avoided as much as possible but instead inject the actual time value or generating time function to functions or objects. This will allow tests to full control of what the time is and write us deterministic code.
Also, injecting dependencies help us to mock the responses of http servers an application sends a request to.
Second, having conditional logics for each environment should not exist, like change a logic based on development, staging, or production environment, or even logic runs only in OS or CI specific condition. Instead, add a new configuration or environment variable which allows test codes to control these logics based on the variables and have consistency of production code across environments.
Third, avoid using global states as much as possible. Global variables are non-negotiable. Other global states like local files, or environment variables should not be used only for sharing data. They should be used because applications or libraries require them.


On top of the above production code requirement, we must write easy to read and write test codes because test code itself can easily become hard to read.
Following simplicity and readability's section is must-have, but test codes must be taken care of from the nature of test code.
Almost all of test code is consists of followings:
1. initialize global states like databases. For example, records on the table must be cleaned up before any tests to avoid getting affected by other tests. For a function with a local file, create a temporary file or a file in a temporary directory, so these files will be automatically deleted. An external API's client must be initialized with mock servers, too.
2. Initialize local variables, objects or states with or without fixture files.
3. Run an actual test code. and
4. Verify the result.

Considering these typical test pattern, there are a few techniques we should follow to ensure test code readability and extensibility.
At first, each test cases should be written with table driven tests. This allows us to write test code share the same logic for testing but just use the different data. The exceptions are if it's very different test case for some steps described in the above steps.
Second, share the initialization logics for global states, objects or mocks. An object defined in an application is generated in a few or numerous test cases, so having a factory function to fill in the default values help to minimize assigning variables on each test case which relates to that test case, and also reduce the changes when there is an change against that object. Similar things for records in a DB table. For a DB record, factory functions for a record type like ORM help but also have insert functions with vararg types to allow any number of record can be inserted easily.
Third, in order to ensure integration testings without mocks, it's important to have an object to initialize many components at once and share it among test cases. For example, many http handlers need to use a DB client, mock client of an external API, and generator for timestamps. In this case, it's much easier to have one object to hold these test objects and use these objects in various test cases for HTTP handlers increase reusability and consistency.

There are a few more things to ensure high testability.
First, it's important not to use mocks as much as possible. This can detect issues by integration testings. RDB server can be runnable easily on a local environment by a container and should not be mocked. The 3rd party's API or other API's response is hard to build locally so they should be mocked during unit or integration testing.

Besides, for AI agents like Claude Code, there are a few important things
First, it's important not to generate embedding base64 or other encoded strings for test data like images or certificates because they are not human readable. Instead, generate image or certificates files and read those files from test codes.
Second, make sure to keep or increase a test coverage. It's easy to skip writing test cases, but AI agents must not skip the test cases.



## Performances

The most significant part of latency comes from I/O, especially network I/Os like requests of external APIs or running SQL queries against a database.
It's more important to make this performance better than simple but too slow code.
In order to achieve high performance, make sure these things.
First, batch multiple items and records into one request. For example, request to a DB database can be batched by one query instead of sending multiple queries to get the same result. This helps to reduce latency of an application.
If a batch processing is not possible, consider to run queries concurrently. For example, some external APIs do not support batch requests and in this case, sending external APIs concurrently helps to shorten the end-to-end time of your process like handling a http request or a batch processing. However, some programming languages are hard to use concurrencies natively. In this case, consider to use producer-subscriber's patterns with a message queue middleware, only when it's possible. If you don't need a response but just want to send a request like sending a notification, then message queue should work pretty well.


## Error handling

The error handling is one of the most important perspective of coding and makes the code complicated.
In the first place, errors must be handled appropriately, and each error should be wrapped to propagate a context information to a call stack. As a rule of thumb, a general error handling must be in-place so all errors are handled appropriately, and on top of that, each error handling must be added to handle each error.

### Server side error handling

First of all, for an application like a HTTP or gRPC server, make sure it returns an appropriate status code. For instance, a client sends a request without an appropriate credential or a request body contains validation error, returns 4XX HTTP response status code. Or if there is an issue happens on a server side while handling a vaild request, reply a server side error like 5XX HTTP response status code.
On top of the above, when a server replies an error response, an error response must include a user friendly message with an action so a user knows what they must do to fix an error. For a server side error, do not include any application internal error message to an end user. Instead, use a generic error message to avoid exposing a security vulnerability to an end user.


### Client side error handling

When an application has to send a request to an external API, this request can fail intentionally or unintentionally.
When an application is designed and need a retry logic, make sure to check an error response and retry errors appropriately.
For example, a request error due to the rate limit, or server side error, or time out error can be retried and in this case, an application should retry following some retry algorithm like exponential backoff.
However, in some systems, an application doesn't need these logics and it's supported on the infrastructure layer.


## Database

Database transaction is one of the most tricky thing.
As a rule of thumb, make sure to create only one DB transaction for writing records in multiple tables for ACID in a request context like handling a HTTP request.
If there are multiple transactions, make sure that they can be inconsistent if there is an error happens between a transaction.
However, it makes sense for a batch processing to have multiple transactions, but in this case, create one DB transaction for update multiple records in multiple tables by one DB transaction, and just repeat this transaction multiple times.

For a database query, always make sure that indexes are used appropriately for SQL queries, especially a query to look up some numbers of records, for example more than 10,000 records.



## Monitoring

In order to make sure an application is reliable and perform well after running the application on production, telemetry data like metrics, logs, or traces are handy.
Open Telemetry is the industry standard and available across multiple languages, so set up its SDK on the bootstrap process and add instrumentation when it's needed.
Telemetry must be exported into other components so that they can be monitored appropriately.

### Metrics

For metrics, export RED metrics for HTTP or gRPC servers at least. But whenever other metrics are required, instrument them to output metrics. But avoid adding some parameters to increase high cardinality because some backend for metrics uses TSDB and it doesn't handle high cardinality well.


### Logs

Logs are helpful to see the details of what's going on in an application, especially in a production. To filter out logs appropriately, use severity appropriately as follows:
- debug: if it's not necessary on production but still need to check during development
- info: if it's important to see the details of some data
- warning: output an expected error case with some error information
- error: output an unexpected error case with error information in order to look into why the error happens later, for example, server side errors of an HTTP server.

Besides, logs must be structured appropriately. And include contexts like a trace span id, a user id of a HTTP request, some of request parameters to help looking into an applications' errors easily. However, make sure not to include sensitive information. Because of this reason, environment variables, an entire body of HTTP requests or responses should not be output in logs, or should be output but redacted appropriately.

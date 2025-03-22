# Project Overview 
This is now my second full-stack application. But with a twist.

The reason why I did this is because our property manager needed assistance collecting fees for building repairs
and I've decided to use the opportunity to create a coding project to keep track of who's paying.

The application follows a three-tier architecture pattern (some would say it's even implementing a microservice approach but I wouldn't agree) with these key components:

* Backend Service in Go: RESTful API for data operations
* Frontend Service (Next.js): User interface with React components
* Database (PostgreSQL): Persistent data storage
* CI/CD Pipeline (Jenkins): Automated testing and deployment
* Observability Stack: Monitoring and logging

For the backend I tried to use modular structure as much as possible.
I have defined plenty of constants to avoid hardcoding values, which is a good practice for maintainability. I gotta mention this, even though it would come later in my documentation, that some of the code edits I did were as a result of *Sonarqube*'s suggestions:


![alt text](assets/sonarqube_test2.png)

Above is a sample screenshot of some of the advice their platform gave me after I implemented its functionality in my Jenkins pipeline. It's pretty neat if you get tired of your preferred AI agent's suggestions. So what I've changed because of Sonar is I constructed HTTP Header Constants, Error Message Constants and a MIME (Multipurpose Internet Mail Extension) Type Constant.


### The main function flow:

* Environment loading with godotenv
* Database connection and setup
* Router configuration
* API endpoint registration
* CORS middleware implementation
* Server startup


### API Endpoints:

* /health: Health check endpoint (The main reason I wanted to add this is because I wanted the backend to be as much Kubernetes-ready as possible)
* /api/repair/users: GET (all users) and POST (create user)
* /api/repair/users/{id}: GET, PUT, DELETE operations for specific users. Far from perfect but will come back and make it more modular after some time.


### Middleware:

* CORS enablement
* JSON content type setting

It's all far from production ready but I wanted to really focus on the infrastructure components later in another repository.

# The Continuos Integration Map of my Application

Below is a Figma diagram I made to illustrate all the bits and pieces that came into play on the surface level. 


![alt text](assets/application_diagram.png)


# It All Started With Jenkins

As I was building my initial version of the pipeline I wanted to implement a secondary branch where the hypothetical dev, ops and QA teams can build and test new features and debug. I also wanted to create for them a nice pipeline to do so. Well, because GitHub can only send webhook payloads to a publicly accessible URL I needed to set-up:

* A fixed local IP via DHCP reservation, so my router always knew where Jenkins lived.
* Port forwarding from my router (ports 80 and 443) to my Jenkins machine.
* DuckDNS domain to point to my home’s dynamic public IP.
* Nginx reverse proxy to serve Jenkins at that domain with SSL.
* Let's Encrypt / Certbot to generate valid HTTPS certs (GitHub requires HTTPS for webhooks).

### Final result: 

![alt text](<assets/Screenshot from 2025-03-21 21-44-40.png>)

### My goal was to create a smooth workflow like:

- Devs push code to GitHub.
- GitHub sends a webhook to Jenkins.
- Jenkins triggers my CI/CD pipeline (test, build, push Docker images, etc.).
- I get feedback in Jenkins and/or GitHub status checks.

My very last version of my Jenkins pipeline looked like this: 

![alt text](<assets/Screenshot from 2025-03-18 19-58-26.png>)

120 Builds is an old number by the way.

- I have implemeneted parameters and environment variables to make the pipeline flexible and reduce hardcoded values.

- I leveraged Jenkins credentials binding for secure access to sensitive data (Docker tokens, GitHub API tokens) to adhere to security best practices.

- I fixed an issue where babel was interfering with the nextjs engine when building the dockerfile so I had to tinker with the babelrc config file and ultimately resolved the problem.

- I managed to set-up and run unit tests concurrently for both frontend and backend to minimize build time and improve efficiency. I had some issues with setting up the Cypress E2E test in headless mode as Jenkins runs are obviously a non GUI environment but I am so proud that I managed to crack it by installing the xvfb dependency which is an in-memory display server for a UNIX-like operating system.
Below is a screenshot of a local cypress run test I did in the beginning:

![alt text](<assets/Screenshot from 2025-02-13 17-45-42.png>)

- The dynamic tagging of Docker images using parameterized version numbers ensures consistent and reproducible builds.This is me adding parameters for the builds if one day I decide to scrap the automatic run on push to feature approach:

![alt text](<assets/Screenshot from 2025-03-18 19-48-19.png>)

- What my pipeline also does is it helps the hypothetical DevOps/Cloud team with infrastructure changes by automatating PR creation by creating a pull request to update the infrastructure repository with new and tested application code. This way I ensure that changes are manually reviewed and merged through a controlled process.

- Okay but why do I need PostgreSQL for SonarQube? Well, it's mainly for persistence and performance. Also, it requires a database in order to store the code analysis results, the quality metrics history, the user accounts and permissions and the project's configurations. Many people might object that I've configured SonarQube to use the same PostgreSQL instance as my application but in my defense SonarQube creates its own schema within my database and it keeps its tables separate from my application data. Also, it's all done in the dev stage of the SDLC so get off my ass. It has nothing to do with the production database. But anyway I've configured persistent volumes for SonarQube's data and logs and also I made sure that Sonarqube always starts after the db container to prevent errors. 

- Dynamic Branch and PR Naming: Functions that generate branch names and PR titles based on updated components definitely improve clarity and traceability.

- And finally I have the integrated code quality checks by running SonarQube as a Docker container, installing sonar-scanner locally and adding sonar-project.properties to my repo and wiring it into Jenkins pipeline.

- At the end I ensure that the workspace is cleaned up after each build and provides clear feedback on the build status.

- I really also wanted to mention that I have integrated into my pipeline a *Sonatype Nexus Repository*

![alt text](<assets/Screenshot from 2025-03-08 12-22-12.png>) 


![alt text](<assets/Screenshot from 2025-03-08 12-21-32.png>)


It is really worth mentioning that my build time got reduced by about 20-25 seconds for such a small application like mine, compared to pushing to docker hub. Imagine that on scale? But because I later implemented ArgoCD it was far easier for me to switch to dockerHub because I think I had to set-up Image Updater but I am not quite sure yet how or if I even needed to do that. Anyway, it's very much worth sharing that with you.


# Who Is Carrying The Load?

My initial plan for the Load and Performance Testing was to set-up a complex JMeter -> influxDB/Prometheus -> Grafana pipeline but at that point I really wanted to move on and start building the infrastructure. I still did two simple stand-alone load tests.

### Parameters 

- I set it to 400 Number of Threads (Users) with a ramp-Up period of 120 seconds.

- The duration is 300 seconds (5 minutes).

- The path being /api/repair/users with a GET Method. The 404 Page is /test.


![alt text](<assets/Screenshot from 2025-01-24 21-35-37.png>) 

This is what a successful test looks like. There were 500 total requests. 250 for the "home page" and 250 for the 404 Page. I configured the assertion for to treat 404 as a successful response. The Average Response Time (ms) is:

- Home Page: 1.60

- Page Returning 404: 1.04

- Max Response Time: 27 ms

- Throughput: ~49.3 transactions/sec overall

![alt text](<assets/Screenshot from 2025-01-24 21-03-45.png>)

Now, in the screenshot above I don't treat the /test path (404 Page) as a success so I get an APDEX Score of 0.000. The slightly higher throughput: ~80 TPS is maybe due to the shorter error-handling time.


# My Observability Set-Up











# Brief Demo



# Final Thoughts

I am definitely missing a lot of things during the development of this project as I think it covered plenty of areas in the SDLC. All in all, it was a very steep learning curve, one that many people would've given up on, but it was definitely worth it in order to understand why being 
Speaking of DevOps this project was just the beginning.


 
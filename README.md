# Overview 
This is my second full-stack application. This time with completely different technologies.
I used Go, Typescript with Next.js, little bit of tailwind, PostgreSQL and my favorite Docker Compose.

The reason why I did this is because our property manager needed assistance collecting fees for building repairs
and I've decided to use the opportunity to create a coding project to keep track of who's paying.

It's far from production ready but I just had to get something up and running as fast as I can. I'll improve on it in the future.
The code that's on github uses the default postgres password but I'm running it on my machine with .env variables.

All three main services are on separate containers and also pushed to docker hub.

![alt text](assets/volume.png)

# Usage

You can clone this repository and then `pull` the docker images from docker hub:

`docker pull tankengine/nextapp:1.0.0`
`docker pull tankengine/goapp:1.0.0`
`docker pull tankengine/postgres:15`

After that you can start the services with `docker-compose up` and go to `http://localhost:3000`.

![alt text](assets/312321321.png)

The above screenshot is a random address with random people added.

# Brief Demo

Alternatively, if you don't want to toy around with docker images and stuff here's a short .gif demonstration.


![alt text](assets/Untitled-ezgif.com-optimize.gif)

And also a VIDEO of me explaining my process:



Problems I fixed:

I mounted the .env file in the docker compose file. 
Then i removed the cached postgres db and restarted the docker images.
Now it was able to reconnect.


Also pushing to docker hub with re-tagging an image.
Useful commands "docker images" and "docker tag d4332ddfd789 tankengine/goapp:latest" and then push with "docker push tankengine/goapp:latest"


Jenkins part:

I've decided to use a local installation of Jenkins instead of having it as container and mindlessly fiddling with docker groups, Docker-IN-Docker images, sockets
and user permissions.
I added the jenkins user to the docker group: sudo usermod -aG docker jenkins and did a test script:

Started by user jenkins
[Pipeline] Start of Pipeline
[Pipeline] node
Running on Jenkins in /var/lib/jenkins/workspace/docker test
[Pipeline] {
[Pipeline] stage
[Pipeline] { (Docker Test)
[Pipeline] script
[Pipeline] {
[Pipeline] sh
+ docker ps
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
[Pipeline] }
[Pipeline] // script
[Pipeline] }
[Pipeline] // stage
[Pipeline] }
[Pipeline] // node
[Pipeline] End of Pipeline
Finished: SUCCESS


fixed an issue where babel was interfering with the nextjs engine when building the docker compose so I had to tinker with the babelrc config file.
Now the jest test works and the compose works as well.

did the jest db, frontend and backend tests. Did the E2E test.

Created credentials for the dockerhub account inside Jenkins.

Set up Cypress in the headless environment in Jenkins.rm -rf cypress.
Installed the xvfb dependency. Config the file for no support file.

I did a DHCP reservation for a local static IP so I can host my jenkins on an nginx server
in order to expose it to GItHub's webhook so I can build on push to the repo.

I forwarded requests on port 80 and port 443 from my public IP to port 80 and 443 on my local Ubuntu server.

Nexus:

I had to create an npmrc file which is a configuration file used by the npm(Node Package Manager) command-line tool. It allows you to customize various settings related to how npm behaves while managing packages and dependencies for your Node.

The commands below fixed the issue of not uploading my npm dependencies to Nexus: 

npm cache clean --force
rm -rf node_modules
npm install


I managed to reduce the build time from more than 2 minutes to a whopping 17 seconds!!!
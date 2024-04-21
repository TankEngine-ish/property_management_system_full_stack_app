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


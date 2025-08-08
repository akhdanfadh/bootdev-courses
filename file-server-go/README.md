# Learning File Servers and CDNs with S3 and CloudFront

Welcome to *large* file storage! Building a (good) web application *almost always* involves handling "large" files of some kind - whether its static images and videos for a marketing site, or user generated content like profile pictures and video uploads, it always seems to come up.

In this course we'll cover strategies for handling files that are kilobytes, megabytes, or even *gigabytes* in size, as opposed to the small structured data that you might store in a traditional database (integers, booleans, and simple strings).

## Learning Goals

- Understand what "large" files are and how they differ from "small" structured data
- Build an app that uses [AWS S3](https://aws.amazon.com/s3/) and [Go](https://www.boot.dev/courses/learn-golang) to store and serve assets
- Learn how to manage files on a "normal" (non-s3) filesystem based application
- Learn how to *store and serve* assets at scale using serverless solutions, like AWS S3
- Learn how to *stream* video and to keep data usage low and improve performance

## AWS Account Required

This course will require an AWS account. We will *not* go outside of the [free tier](https://aws.amazon.com/free/), so if you do everything properly you shouldn't be charged. That said, you will need to have a credit card on file, and if you do something wrong you *could* be charged, so just be careful and understand the risk.

We recommend *deleting all the resources* that you create when you're done with the course to avoid any charges. We'll remind you at the end.

## Tubely

In this course we'll be building "Tubely", a SaaS product that helps YouTubers manage their video assets. It allows users to upload, store, serve, add metadata to, and version their video files. It will also allow them to manage thumbnails, titles, and other video metadata.

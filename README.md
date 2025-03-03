# **URL Shortener - Project Documentation**  

## **Table of Contents**  
1. [Introduction](#introduction)  
2. [Project Approach](#project-approach)  
   - [Phase 1: Core URL Shortening and Redirect](#phase-1-core-url-shortening-and-redirect)  
   - [Phase 2: Enhancements - Validations & Expiration](#phase-2-enhancements---validations--expiration)  
   - [Phase 3: MongoDB Integration & Usage Tracking](#phase-3-mongodb-integration--usage-tracking)  
3. [Design Decisions](#design-decisions)  
   - [Architecture: MVC Pattern](#architecture-mvc-pattern)  
   - [Database Choice: MongoDB over PostgreSQL](#database-choice-mongodb-over-postgresql)  
   - [Future Scalability Considerations](#future-scalability-considerations)  
4. [Challenges Faced](#challenges-faced)  

---

## **1. Introduction and To-Run Instructions**  
The **URL Shortener** project is a backend service that takes long URLs and generates short, unique identifiers that can be used to access the original link. It also tracks URL usage statistics and handles expiration for URLs.  

The project follows a **structured MVC architecture**, uses **MongoDB** for storage, and is built using **Golang with the Gin framework**.  

To Run the project, clone the repo, then paste your own MONGO_URI string in it. Then in the root project directory, use the command `go run main.go`

---

## **2. Project Approach**  

The project was built in three structured phases, each focusing on incremental improvements and new feature additions.

### **Phase 1: Core URL Shortening and Redirect**  
#### **Goal:**  
To build the fundamental functionality of shortening a URL and redirecting a user when they access the short link.

#### **Implementation Details:**  
1. **Project Structure (MVC Architecture)**  
   - **Controllers** handle API requests and responses.  
   - **Services** handle business logic such as generating short URLs.  
   - **Storage Layer** manages persistent storage of URLs in MongoDB.  

2. **Shortening a URL**  
   - A request is made to `/shorten` with a long URL.  
   - The system generates a unique short code.  
   - The short code is stored along with the original URL in MongoDB.  

3. **Redirecting a URL**  
   - A request is made to `/short/{shortURL}`.  
   - The system fetches the long URL from storage and redirects the user.  

### **Phase 2: Enhancements - Validations & Expiration**  
#### **Goal:**  
Improve reliability by adding URL validation and expiration logic.

#### **Implementation Details:**  
1. **URL Validation:**  
   - Ensures only valid URLs with `http` or `https` are accepted.  
   - Uses `net/url` package to parse and validate URLs.  

2. **Expiration Logic:**  
   - Each short URL is assigned an `expiration_time`.  
   - If a URL is accessed after its expiration, it is deleted from the database.  

### **Phase 3: MongoDB Integration & Usage Tracking**  
#### **Goal:**  
Persist data in MongoDB and track how many times each short URL is accessed.

#### **Implementation Details:**  
1. **MongoDB as Primary Storage:**  
   - Implemented `mongodb_store.go` for MongoDB interactions.  
   - Used MongoDB's BSON format to store URLs, expiration time, and usage count.  

2. **Usage Tracking:**  
   - Each short URL has a `usage` field initialized to `0`.  
   - Every time the URL is accessed, `usage` is incremented.  
   - A new API `/usage/{shortURL}` was implemented to retrieve usage statistics.  

---

## **3. Design Decisions**  

### **Architecture: MVC Pattern**  
I chose the **Model-View-Controller (MVC)** pattern because:  
- It keeps concerns separate, making the code more maintainable.  
- It scales well for small to mid-sized applications.  
- Other options like **Microservices Architecture or Clean Architecture** were unnecessary for this project since it's a relatively small-scale application.

### **Database Choice: MongoDB over PostgreSQL**  
- **MongoDB** was chosen because:  
   - It provides **flexibility** to add new fields dynamically (e.g., tracking analytics).  
   - The schema-less nature is suitable for a project where data models may evolve.  
   - Queries for URL lookups are **fast**, even with large amounts of data.

- **PostgreSQL (SQL-based DB)** was considered but rejected because:  
   - Schema constraints make it **less flexible** for future modifications.  
   - For read-heavy operations like URL lookups, NoSQL databases offer **better performance**.  

### **Future Scalability Considerations**  
If this application needs to scale further:  
1. **Redis for Caching:**  
   - Using **Redis** as a cache layer would significantly improve lookup performance for high traffic.  
   - URLs could be stored in Redis with TTL (Time-To-Live) to reduce database queries.  

2. **Load Balancing & Rate Limiting:**  
   - Adding **NGINX** or **API Gateway** would help distribute traffic efficiently.  
   - Rate limiting can be enforced to prevent abuse (e.g., **throttling API requests**).  

---

## **4. Challenges Faced**  

### **1. Designing Expiration Logic for Short URLs**  
🔹 **Problem:** URLs need an expiration mechanism, but MongoDB does not provide automatic deletion for expired records.  
🔹 **Solution:** Implemented manual expiration checking before serving a request and deleting expired URLs on access.  

### **2. Understanding URL Shortener System Design**  
🔹 **Problem:** Needed to learn how real-world URL shorteners like **Bitly** or **TinyURL** work.  
🔹 **Solution:** Researched common system design patterns for URL shorteners and applied relevant concepts.  

### **3. SQL vs NoSQL - Choosing the Right Database**  
🔹 **Problem:** Should I use PostgreSQL or MongoDB?  
🔹 **Solution:** Evaluated the pros and cons of each and decided on MongoDB for flexibility and fast lookups.  

### **4. Caching Considerations**  
🔹 **Problem:** Should I integrate Redis for quick lookups?  
🔹 **Solution:** Since the assignment does not require large-scale optimization, I skipped Redis. But I now understand the circumstances under which caching would be beneficial.  



🎯 **Project Completed Successfully!** It was a great journey getting to build this application, hehe! 🚀

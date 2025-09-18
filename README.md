# 🏛️ GoForum

GoForum is a modern, feature-rich forum platform built with Go and designed for sharing ideas, discussions, and community building. It provides a clean, user-friendly interface for both users and administrators.

## 🌟 Features

### For Users
- 📝 **Create & Manage Posts** - Write forum posts with rich text editor
- 💾 **Save Favorite Posts** - Bookmark posts for quick access
- 💬 **Comment System** - Engage with nested comments and replies
- ❤️ **Like Comments** - Show appreciation for valuable contributions
- 👤 **User Profiles** - Manage personal content and view statistics
- 🏷️ **Category Browsing** - Explore posts by topics

### For Administrators
- 🛠️ **Admin Dashboard** - Complete content management system
- 📊 **User Management** - Monitor and manage community members
- 🏷️ **Category Management** - Create and organize content categories
- ✅ **Content Moderation** - Approve/reject posts and comments
- 📈 **Analytics** - Track engagement and community growth

## 🛠️ Tech Stack

- **Backend**: Go 1.25+
- **Web Framework**: Custom HTTP router (julienschmidt/httprouter)
- **Database**: SQLite with GORM
- **Session Management**: Gorilla Sessions
- **Frontend**: HTML Templates, Bootstrap, JavaScript
- **Rich Text Editor**: Summernote
- **Icons**: Font Awesome

## 📦 Installation

### Prerequisites

- Go 1.25 or higher
- Git

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/sulecengiz/goforum.git
   cd goforum
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

4. **Access the application**
   - Frontend: http://localhost:8080
   - Admin Panel: http://localhost:8080/admin

### Default Admin Credentials
- Username: `admin`
- Password: `admin123` (Please change after first login)

## 📁 Project Structure

```
goforum/
├── admin/              # Admin panel
│   ├── controllers/    # Admin controllers
│   ├── models/         # Admin data models
│   ├── views/          # Admin templates
│   └── assets/         # Admin static files
├── site/               # Public site
│   ├── controllers/    # Site controllers
│   ├── models/         # Site data models
│   ├── views/          # Site templates
│   └── assets/         # Site static files
├── config/             # Configuration
├── uploads/            # User uploaded files
└── main.go             # Application entry point
```

## 🚀 Key Features Breakdown

### Forum Management
- Create, edit, and delete forum posts
- Rich text editor with image upload support
- Category-based organization
- Post approval workflow

### User Interaction
- User registration and authentication
- Profile management
- Comment system with nested replies
- Like/dislike functionality
- Save posts for later reading

### Admin Features
- Complete dashboard for content management
- User role management
- Content moderation tools
- Category management
- System analytics

## 🎨 Screenshots

<img width="1905" height="944" alt="image" src="https://github.com/user-attachments/assets/fb5a93b9-3780-4b21-8877-65d0583b5173" />
<img width="1905" height="942" alt="image" src="https://github.com/user-attachments/assets/f4d5e9fa-cdb8-4174-8cbf-70745fcfd02a" />
<img width="1917" height="940" alt="image" src="https://github.com/user-attachments/assets/ac24db20-3b04-4eb6-9826-31917ab64c96" />
<img width="1903" height="943" alt="image" src="https://github.com/user-attachments/assets/90f72fe0-5620-4c02-bbb2-84573147589e" />
<img width="1906" height="942" alt="image" src="https://github.com/user-attachments/assets/26e811f6-28a3-4a23-9a25-6deec58576f9" />


<img width="1917" height="942" alt="image" src="https://github.com/user-attachments/assets/d002c43c-96ec-48a1-9f31-ef66454eda42" />
<img width="1912" height="944" alt="image" src="https://github.com/user-attachments/assets/2c4e23bd-d43b-44a7-b905-55915a24d512" />
<img width="1917" height="941" alt="image" src="https://github.com/user-attachments/assets/dc727134-15ce-4d6d-a780-75cb0aeca7e1" />
<img width="1918" height="940" alt="image" src="https://github.com/user-attachments/assets/13960010-295f-4d60-8a93-9055649f5980" />
<img width="1904" height="946" alt="image" src="https://github.com/user-attachments/assets/c5f4c906-6545-4187-8ba9-fa5153512fcd" />






## 🔧 Configuration

The application uses SQLite by default and creates the database automatically on first run. All configurations are handled through environment variables and the main configuration file.

## 📝 API Endpoints

### Public Routes
- `GET /` - Homepage
- `GET /yazilar/:slug` - Post details
- `POST /yazilar/:slug/yorum-ekle` - Add comment
- `GET /profile` - User profile
- `POST /save-post/:postID` - Save/unsave post

### Admin Routes
- `GET /admin` - Admin dashboard
- `POST /admin/add` - Create new post
- `GET /admin/edit/:id` - Edit post
- `DELETE /admin/delete/:id` - Delete post

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Commit your changes**
   ```bash
   git commit -am 'Add some feature'
   ```
5. **Push to the branch**
   ```bash
   git push origin feature/your-feature-name
   ```
6. **Create a Pull Request**

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Süleyman Cengiz**
- GitHub: [@sulecengiz](https://github.com/sulecengiz)
- LinkedIn: [sulecengizz](https://www.linkedin.com/in/sulecengizz/)

## 🙏 Acknowledgments

- Built with ❤️ using Go
- Thanks to the open-source community for the amazing tools and libraries
- Special thanks to all contributors

---

**GoForum** - Building communities, one post at a time! 🌟

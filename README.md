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

*Add screenshots of your application here*

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

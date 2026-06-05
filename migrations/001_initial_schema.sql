-- Users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Books
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    year INTEGER,
    available_copies INTEGER DEFAULT 0,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Borrows
CREATE TABLE IF NOT EXISTS borrows (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    borrow_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Favorites
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, book_id)
);

-- Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
CREATE INDEX IF NOT EXISTS idx_books_category ON books(category);
CREATE INDEX IF NOT EXISTS idx_borrows_user ON borrows(user_id);
CREATE INDEX IF NOT EXISTS idx_borrows_status ON borrows(status);
CREATE INDEX IF NOT EXISTS idx_favorites_user ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);

-- Seed data
INSERT INTO books (title, author, description, category, year, available_copies, image_url) VALUES
('The Great Gatsby', 'F. Scott Fitzgerald', 'A classic American novel', 'Fiction', 1925, 5, 'https://covers.openlibrary.org/b/id/7222246-M.jpg'),
('To Kill a Mockingbird', 'Harper Lee', 'Southern Gothic novel', 'Fiction', 1960, 3, 'https://covers.openlibrary.org/b/id/8228691-M.jpg'),
('1984', 'George Orwell', 'Dystopian social science fiction', 'Science Fiction', 1949, 7, 'https://covers.openlibrary.org/b/id/7222246-M.jpg'),
('Pride and Prejudice', 'Jane Austen', 'Romantic novel', 'Romance', 1813, 4, 'https://covers.openlibrary.org/b/id/8228691-M.jpg'),
('The Hobbit', 'J.R.R. Tolkien', 'Fantasy adventure', 'Fantasy', 1937, 6, 'https://covers.openlibrary.org/b/id/7222246-M.jpg')
ON CONFLICT DO NOTHING;
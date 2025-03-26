CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    recipient VARCHAR(20) NOT NULL,
    content TEXT NOT NULL CHECK (CHAR_LENGTH(content) <= 255),
    status ENUM('pending', 'processing', 'sent', 'failed') DEFAULT 'pending',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
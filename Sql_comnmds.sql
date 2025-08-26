USE booksdb;
CREATE TABLE books (
  id VARCHAR(36) PRIMARY KEY,
  title VARCHAR(255),
  author VARCHAR(255),
  quantity INT
);

INSERT INTO books (id, title, author, quantity) VALUES
('1', 'In Search of Lost Time', 'Marcel Proust', 2),
('2', 'The Great Gatsby', 'F. Scott Fitzgerald', 5),
('3', 'War and Peace', 'Leo Tolstoy', 6);

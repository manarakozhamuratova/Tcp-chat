CREATE TABLE IF NOT EXISTS users(
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts(
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (user_id)
);


CREATE TABLE IF NOT EXISTS categories(
	category_id INTEGER PRIMARY KEY AUTOINCREMENT,
	category TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS post_categories (
		post_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (post_id, category_id),
		FOREIGN KEY (post_id) REFERENCES posts (post_id),
		FOREIGN KEY (category_id) REFERENCES categories (category_id)
);

CREATE TABLE IF NOT EXISTS sessions(
	session_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	token TEXT NOT NULL,
	expiry DATE NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id)	ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS comments(
	comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	message TEXT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id),
	FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);
	

CREATE TABLE IF NOT EXISTS posts_likes_dislikes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	like INTEGER NOT NULL,
	dislike INTEGER NOT NULL,
	FOREIGN KEY (post_id) REFERENCES posts (post_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments_likes_dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		like INTEGER NOT NULL,
		dislike INTEGER NOT NULL,
		FOREIGN KEY (comment_id) REFERENCES comments (comment_id) ON DELETE CASCADE
);

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Adventure stories' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Adventure stories') LIMIT 1;

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Crime' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Crime') LIMIT 1;

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Fantasy' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Fantasy') LIMIT 1;

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Humour' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Humour') LIMIT 1;

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Mystery' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Mystery') LIMIT 1;

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Plays' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Plays') LIMIT 1;

INSERT INTO categories (category)
	SELECT * 
	FROM (SELECT 'Other' as category) AS tmp 
	WHERE NOT EXISTS (SELECT category FROM categories WHERE category='Other') LIMIT 1;

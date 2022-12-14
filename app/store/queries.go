package store

const (
	createProduct = `
	INSERT INTO product (product_name)
	VALUES ($1) RETURNING product_id;`

	createProductArticle = `
	INSERT INTO product_article (product_id, article_id, article_amount)
	VALUES ($1, $2, $3);`

	createArticle = `
	INSERT INTO article (article_id, stock, article_name)
	VALUES ($1, $2, $3);`

	getProductArticlesByProductID = `
	SELECT article_id, article_amount FROM product_article
	WHERE product_id = $1`

	updateArticleStockForSellProduct = `
	UPDATE article SET stock = stock - $1
	WHERE article_id = $2;`

	getProductsWithStock = `
	SELECT product_article.product_id, MIN(article.stock / product_article.article_amount) as stock FROM product_article
	LEFT JOIN article ON article.article_id = product_article.article_id
	GROUP BY product_article.product_id;`
)

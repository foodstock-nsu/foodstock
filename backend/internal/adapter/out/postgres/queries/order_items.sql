-- name: CreateOrderItem :exec
INSERT INTO order_items (
    id,
    order_id,
    item_id,
    item_amount,
    price_at_purchase
) VALUES (
    @id,
    @order_id,
    @item_id,
    @item_amount,
    @price_at_purchase
);

-- name: CreateOrderItemsBatch :copyfrom
INSERT INTO order_items (
    id,
    order_id,
    item_id,
    item_amount,
    price_at_purchase
) VALUES (
    @id,
    @order_id,
    @item_id,
    @item_amount,
    @price_at_purchase
);

-- name: ListOrderItems :many
SELECT
    id,
    order_id,
    item_id,
    item_amount,
    price_at_purchase
FROM order_items
WHERE order_id = @order_id;
class Order {
  final String id;
  final String userId;
  final String contractId;
  final String type;
  final String status;
  final int quantity;
  final int quantityFilled;
  final double? limitPrice;
  final DateTime createdAt;

  Order({
    required this.id,
    required this.userId,
    required this.contractId,
    required this.type,
    required this.status,
    required this.quantity,
    required this.quantityFilled,
    this.limitPrice,
    required this.createdAt,
  });

  factory Order.fromJson(Map<String, dynamic> json) {
    return Order(
      id: json['id'],
      userId: json['user_id'],
      contractId: json['contract_id'],
      type: json['type'],
      status: json['status'],
      quantity: json['quantity'],
      quantityFilled: json['quantity_filled'] ?? 0,
      limitPrice: json['limit_price_credits'] != null
          ? (json['limit_price_credits']).toDouble()
          : null,
      createdAt: DateTime.parse(json['created_at']),
    );
  }

  bool get isFilled => status == 'FILLED';
  bool get isActive => status == 'ACTIVE' || status == 'PARTIALLY_FILLED';
}

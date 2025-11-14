class Market {
  final String id;
  final String question;
  final String description;
  final String category;
  final String status;
  final DateTime resolutionDate;
  final DateTime createdAt;
  final String? outcome;
  final List<Contract>? contracts;

  Market({
    required this.id,
    required this.question,
    required this.description,
    required this.category,
    required this.status,
    required this.resolutionDate,
    required this.createdAt,
    this.outcome,
    this.contracts,
  });

  factory Market.fromJson(Map<String, dynamic> json) {
    return Market(
      id: json['id'],
      question: json['question'],
      description: json['description'] ?? '',
      category: json['category'],
      status: json['status'],
      resolutionDate: DateTime.parse(json['resolution_date']),
      createdAt: DateTime.parse(json['created_at']),
      outcome: json['outcome'],
      contracts: json['contracts'] != null
          ? (json['contracts'] as List).map((c) => Contract.fromJson(c)).toList()
          : null,
    );
  }
}

class Contract {
  final String id;
  final String marketId;
  final String side;
  final double currentPrice;
  final int volume;

  Contract({
    required this.id,
    required this.marketId,
    required this.side,
    required this.currentPrice,
    required this.volume,
  });

  factory Contract.fromJson(Map<String, dynamic> json) {
    return Contract(
      id: json['id'],
      marketId: json['market_id'],
      side: json['side'],
      currentPrice: (json['current_price_credits'] ?? 0).toDouble(),
      volume: json['volume'] ?? 0,
    );
  }
}

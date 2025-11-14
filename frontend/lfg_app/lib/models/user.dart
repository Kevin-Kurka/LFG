class User {
  final String id;
  final String email;
  final String displayName;
  final String status;
  final DateTime createdAt;

  User({
    required this.id,
    required this.email,
    required this.displayName,
    required this.status,
    required this.createdAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      email: json['email'],
      displayName: json['display_name'] ?? '',
      status: json['status'] ?? 'ACTIVE',
      createdAt: DateTime.parse(json['created_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'display_name': displayName,
      'status': status,
      'created_at': createdAt.toIso8601String(),
    };
  }
}

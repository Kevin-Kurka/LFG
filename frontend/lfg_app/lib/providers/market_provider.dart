import 'package:flutter/foundation.dart';
import 'package:lfg_app/models/market.dart';
import 'package:lfg_app/services/api_client.dart';

class MarketProvider with ChangeNotifier {
  final String? token;
  List<Market> _markets = [];
  Market? _selectedMarket;
  bool _isLoading = false;
  String? _error;
  String? _statusFilter;
  String? _searchQuery;

  List<Market> get markets => _markets;
  Market? get selectedMarket => _selectedMarket;
  bool get isLoading => _isLoading;
  String? get error => _error;

  MarketProvider(this.token);

  Future<void> fetchMarkets({String? status, String? search}) async {
    _isLoading = true;
    _error = null;
    _statusFilter = status;
    _searchQuery = search;
    notifyListeners();

    try {
      final apiClient = ApiClient(token: token);
      _markets = await apiClient.getMarkets(
        status: status,
        search: search,
      );
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> fetchMarket(String marketId) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final apiClient = ApiClient(token: token);
      _selectedMarket = await apiClient.getMarket(marketId);
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
    }
  }

  void clearSelectedMarket() {
    _selectedMarket = null;
    notifyListeners();
  }

  Future<void> refresh() async {
    await fetchMarkets(status: _statusFilter, search: _searchQuery);
  }
}

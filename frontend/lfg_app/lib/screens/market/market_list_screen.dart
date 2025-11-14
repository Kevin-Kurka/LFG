import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:lfg_app/providers/market_provider.dart';
import 'package:lfg_app/screens/market/market_detail_screen.dart';
import 'package:intl/intl.dart';

class MarketListScreen extends StatefulWidget {
  const MarketListScreen({super.key});

  @override
  State<MarketListScreen> createState() => _MarketListScreenState();
}

class _MarketListScreenState extends State<MarketListScreen> {
  final _searchController = TextEditingController();
  String? _statusFilter;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<MarketProvider>(context, listen: false).fetchMarkets();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _handleSearch() async {
    final provider = Provider.of<MarketProvider>(context, listen: false);
    await provider.fetchMarkets(
      status: _statusFilter,
      search: _searchController.text.trim().isEmpty ? null : _searchController.text.trim(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Markets'),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(120),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search markets...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: IconButton(
                      icon: const Icon(Icons.clear),
                      onPressed: () {
                        _searchController.clear();
                        _handleSearch();
                      },
                    ),
                    border: const OutlineInputBorder(),
                    filled: true,
                  ),
                  onSubmitted: (_) => _handleSearch(),
                ),
                const SizedBox(height: 10),
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: [
                      FilterChip(
                        label: const Text('All'),
                        selected: _statusFilter == null,
                        onSelected: (selected) {
                          setState(() {
                            _statusFilter = null;
                          });
                          _handleSearch();
                        },
                      ),
                      const SizedBox(width: 8),
                      FilterChip(
                        label: const Text('Open'),
                        selected: _statusFilter == 'OPEN',
                        onSelected: (selected) {
                          setState(() {
                            _statusFilter = selected ? 'OPEN' : null;
                          });
                          _handleSearch();
                        },
                      ),
                      const SizedBox(width: 8),
                      FilterChip(
                        label: const Text('Resolved'),
                        selected: _statusFilter == 'RESOLVED',
                        onSelected: (selected) {
                          setState(() {
                            _statusFilter = selected ? 'RESOLVED' : null;
                          });
                          _handleSearch();
                        },
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
      body: Consumer<MarketProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.markets.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.error != null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline, size: 48, color: Colors.red),
                  const SizedBox(height: 16),
                  Text('Error: ${provider.error}'),
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: () => provider.fetchMarkets(),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          }

          if (provider.markets.isEmpty) {
            return const Center(
              child: Text('No markets found'),
            );
          }

          return RefreshIndicator(
            onRefresh: () => provider.refresh(),
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: provider.markets.length,
              itemBuilder: (context, index) {
                final market = provider.markets[index];
                return Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  child: InkWell(
                    onTap: () {
                      Navigator.of(context).push(
                        MaterialPageRoute(
                          builder: (context) => MarketDetailScreen(marketId: market.id),
                        ),
                      );
                    },
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Chip(
                                label: Text(market.category),
                                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                              ),
                              const Spacer(),
                              Chip(
                                label: Text(market.status),
                                backgroundColor: market.status == 'OPEN'
                                    ? Colors.green.withOpacity(0.2)
                                    : Colors.grey.withOpacity(0.2),
                                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                              ),
                            ],
                          ),
                          const SizedBox(height: 12),
                          Text(
                            market.question,
                            style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.bold,
                                ),
                          ),
                          if (market.description.isNotEmpty) ...[
                            const SizedBox(height: 8),
                            Text(
                              market.description,
                              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                    color: Colors.grey[600],
                                  ),
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                            ),
                          ],
                          const SizedBox(height: 12),
                          Row(
                            children: [
                              const Icon(Icons.calendar_today, size: 16, color: Colors.grey),
                              const SizedBox(width: 4),
                              Text(
                                'Resolves: ${DateFormat.yMMMd().format(market.resolutionDate)}',
                                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                      color: Colors.grey,
                                    ),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ),
                );
              },
            ),
          );
        },
      ),
    );
  }
}

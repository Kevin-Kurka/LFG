import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:lfg_app/providers/auth_provider.dart';
import 'package:lfg_app/providers/market_provider.dart';
import 'package:lfg_app/providers/wallet_provider.dart';
import 'package:lfg_app/providers/order_provider.dart';
import 'package:lfg_app/screens/splash_screen.dart';
import 'package:lfg_app/screens/auth/login_screen.dart';
import 'package:lfg_app/screens/home/home_screen.dart';

void main() {
  runApp(const LFGApp());
}

class LFGApp extends StatelessWidget {
  const LFGApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProxyProvider<AuthProvider, MarketProvider>(
          create: (_) => MarketProvider(null),
          update: (_, auth, __) => MarketProvider(auth.token),
        ),
        ChangeNotifierProxyProvider<AuthProvider, WalletProvider>(
          create: (_) => WalletProvider(null),
          update: (_, auth, __) => WalletProvider(auth.token),
        ),
        ChangeNotifierProxyProvider<AuthProvider, OrderProvider>(
          create: (_) => OrderProvider(null),
          update: (_, auth, __) => OrderProvider(auth.token),
        ),
      ],
      child: MaterialApp(
        title: 'LFG - Prediction Markets',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(
            seedColor: const Color(0xFF6C5CE7),
            brightness: Brightness.light,
          ),
          useMaterial3: true,
          appBarTheme: const AppBarTheme(
            centerTitle: true,
            elevation: 0,
          ),
        ),
        darkTheme: ThemeData(
          colorScheme: ColorScheme.fromSeed(
            seedColor: const Color(0xFF6C5CE7),
            brightness: Brightness.dark,
          ),
          useMaterial3: true,
          appBarTheme: const AppBarTheme(
            centerTitle: true,
            elevation: 0,
          ),
        ),
        themeMode: ThemeMode.system,
        home: const SplashScreen(),
        routes: {
          '/login': (context) => const LoginScreen(),
          '/home': (context) => const HomeScreen(),
        },
      ),
    );
  }
}
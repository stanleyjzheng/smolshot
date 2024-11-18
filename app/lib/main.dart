import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:platform_device_id/platform_device_id.dart';
import 'package:url_launcher/url_launcher.dart';
import './candles/types.dart';
import './candles/chart.dart';

void main() {
  runApp(SmolShot());
}

// Root of the application
class SmolShot extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Solana App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: PrivateKeyPage(),
    );
  }
}

// Page to input Solana private key
class PrivateKeyPage extends StatefulWidget {
  @override
  _PrivateKeyPageState createState() => _PrivateKeyPageState();
}

class _PrivateKeyPageState extends State<PrivateKeyPage> {
  final TextEditingController _privateKeyController = TextEditingController();

  void _submitPrivateKey() async {
    String privateKey = _privateKeyController.text.trim();
    if (privateKey.isNotEmpty) {
      String? deviceId = await PlatformDeviceId.getDeviceId;

      final response = await http.post(
        Uri.parse('http://localhost:8080/api/v3/set_private_key'),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
        },
        body: jsonEncode(<String, String>{
          'user_id': deviceId ?? 'unknown_device',
          'private_key': privateKey,
        }),
      );

      if (response.statusCode == 200) {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => HomePage(privateKey: privateKey),
          ),
        );
      } else {
        // Show an alert if the request failed
        showDialog(
          context: context,
          builder: (_) => AlertDialog(
            title: Text('Error'),
            content: Text('Failed to save private key. Please try again.'),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: Text('OK'),
              ),
            ],
          ),
        );
      }
    } else {
      // Show an alert if the private key is empty
      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: Text('Error'),
          content: Text('Please enter your private key.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('OK'),
            ),
          ],
        ),
      );
    }
  }

  @override
  void dispose() {
    _privateKeyController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text('Enter Solana Private Key'),
        ),
        body: Center(
          child: Padding(
            padding: EdgeInsets.all(16.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                TextField(
                  controller: _privateKeyController,
                  decoration: InputDecoration(
                    labelText: 'Private Key',
                    border: OutlineInputBorder(),
                  ),
                  obscureText: true, // Hide the input for security
                ),
                SizedBox(height: 20),
                ElevatedButton(
                  onPressed: _submitPrivateKey,
                  child: Text('Submit'),
                ),
              ],
            ),
          ),
        ));
  }
}

// Home page after entering the private key
class HomePage extends StatefulWidget {
  final String privateKey;

  HomePage({required this.privateKey});

  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  double _solBalance = 0.0;
  double _pnutBalance = 0.0;
  final TextEditingController _amountController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _fetchBalances();
  }

  // Fetch user's balances from your API
  void _fetchBalances() async {
    String? deviceId = await PlatformDeviceId.getDeviceId;

    if (deviceId != null) {
      // Fetch SOL balance
      final solResponse = await http.get(
        Uri.parse(
            'http://localhost:8080/api/v3/get_sol_balance?user_id=$deviceId'),
      );

      if (solResponse.statusCode == 200) {
        final solData = jsonDecode(solResponse.body);
        setState(() {
          _solBalance = double.parse(solData['balance'].split(' ')[0]);
        });
      }

      // Fetch PNUT balance
      final pnutResponse = await http.get(
        Uri.parse(
            'http://localhost:8080/api/v3/get_token_balance?user_id=$deviceId&mint_address=7jYfnjn3jHWmUhhgfkZ9uUEKroH3Zz1BvF8RmVQHm86D'),
      );

      if (pnutResponse.statusCode == 200) {
        final pnutData = jsonDecode(pnutResponse.body);
        setState(() {
          _pnutBalance = double.parse(pnutData['balance']) /
              1000000; // Assuming PNUT balance is in micro units
        });
      }
    }
  }

  Future<void> _performTransaction(String action) async {
    String? deviceId = await PlatformDeviceId.getDeviceId;
    String amount = _amountController.text.trim();

    if (amount.isEmpty) {
      // Show an alert if the amount is empty
      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: Text('Error'),
          content: Text('Please enter an amount.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('OK'),
            ),
          ],
        ),
      );
      return;
    }

    final response = await http.post(
      Uri.parse('http://localhost:8080/api/v3/swap_token'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode(<String, dynamic>{
        'user_id': deviceId,
        'input_mint': 'So11111111111111111111111111111111111111112',
        'output_mint': '7jYfnjn3jHWmUhhgfkZ9uUEKroH3Zz1BvF8RmVQHm86D',
        'amount': (double.parse(amount) * 1000000).toInt(),
        'slippage_bps': 100,
      }),
    );

    if (response.statusCode == 200) {
      final responseData = jsonDecode(response.body);
      final transactionSignature = responseData['transaction_signature'];
      final url = 'https://solscan.io/tx/$transactionSignature';

      // Show a dialog with the transaction link
      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: Text('Transaction Successful'),
          content: Text('Transaction Signature: $transactionSignature'),
          actions: [
            TextButton(
              onPressed: () {
                Navigator.pop(context);
                // Open the transaction link
                launch(url);
              },
              child: Text('View on Solscan'),
            ),
          ],
        ),
      );
    } else {
      // Show an alert if the transaction failed
      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: Text('Error'),
          content: Text('Transaction failed. Please try again.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('OK'),
            ),
          ],
        ),
      );
    }
  }

  void _buyCoin() {
    _performTransaction('buy');
  }

  void _sellCoin() {
    _performTransaction('sell');
  }

  @override
  Widget build(BuildContext context) {
    // Example coin data
    Coin coinData = Coin(
      id: 'AS98jKn9RdGYn3MzAiTmUDDPcoWCjdhDfYgwsX8dppP2',
      image: 'https://s2.coinmarketcap.com/static/img/coins/64x64/1.png',
      name: 'Bitcoin',
      shortName: 'BTC',
      price: '123456',
      lastPrice: '123456',
      percentage: '-0.5',
      symbol: 'pnut',
      highDay: '567',
      lowDay: '12',
      decimalCurrency: 18,
    );

    return Scaffold(
      appBar: AppBar(
        title: Text('\$PNUT'),
      ),
      body: Padding(
        padding: EdgeInsets.all(16.0),
        child: Column(
          children: [
            // Display user's SOL balance
            Text(
              'SOL Balance: $_solBalance SOL',
              style: TextStyle(fontSize: 24),
            ),
            SizedBox(height: 10),
            // Display user's PNUT balance
            Text(
              'PNUT Balance: $_pnutBalance PNUT',
              style: TextStyle(fontSize: 24),
            ),
            SizedBox(height: 20),
            // Candlestick chart
            Expanded(
              child: CandleChart(
                coinData: coinData,
                intervalSelectedTextColor: Colors.red,
                intervalTextSize: 20,
                intervalUnselectedTextColor: Colors.black,
              ),
            ),
            SizedBox(height: 20),
            // Amount input field
            TextField(
              controller: _amountController,
              decoration: InputDecoration(
                labelText: 'Amount (SOL)',
                border: OutlineInputBorder(),
              ),
              keyboardType: TextInputType.number,
            ),
            SizedBox(height: 20),
            // Buy and Sell buttons
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                ElevatedButton(
                  onPressed: _buyCoin,
                  child: Text('Buy'),
                ),
                ElevatedButton(
                  onPressed: _sellCoin,
                  child: Text('Sell'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

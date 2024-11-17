import 'dart:convert';

import 'package:http/http.dart' as http;

import './types.dart';

/// fetch candles using api
Future<List<Candle>> fetchCandles(
    {required String symbol, required String interval}) async {
  final uri = Uri.parse(
      "http://localhost:8080/api/v3/coins/$symbol/ohlc?interval=$interval");
  final res = await http.get(uri);

  if (res.statusCode == 200) {
    List<dynamic> data = jsonDecode(res.body);

    /// return candles
    return data.map((e) => Candle.fromJson(e)).toList();
  } else {
    throw Exception('Failed to load candles');
  }
}

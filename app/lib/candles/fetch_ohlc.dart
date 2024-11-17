import 'dart:convert';
import 'dart:math';

import 'package:http/http.dart' as http;

import './types.dart';

Map<String, dynamic> splitTimeUnitWithFullName(String input) {
  // Regular expression to split number and unit
  final regex = RegExp(r'^(\d+)([a-zA-Z]+)$');
  final match = regex.firstMatch(input);

  if (match == null) {
    throw ArgumentError('Invalid input format: $input');
  }

  final number = int.parse(match.group(1)!);
  final unit = match.group(2)!;

  // Mapping shorthand units to full names
  final unitMap = {
    'm': 'minute',
    'h': 'hour',
    'd': 'day',
  };

  if (!unitMap.containsKey(unit)) {
    throw ArgumentError('Unknown unit: $unit');
  }

  final fullUnit = unitMap[unit]!;

  return {
    'number': number,
    'unit': fullUnit,
  };
}

double roundToDecimals(double value, int decimals) {
  final factor = pow(10, decimals);
  return (value * factor).round() / factor;
}

/// fetch candles using api
Future<List<Candle>> fetchCandles(
    {required String symbol,
    required String interval,
    int decimals = 5}) async {
  // parse the period and number of periods
  // e.g. 1d -> 1 day, 1w -> 1 week, 1m -> 1 month

  final timeUnit = splitTimeUnitWithFullName(interval);

  final uri = Uri.parse(
      "http://localhost:8080/api/v3/coins/$symbol/ohlc?period=${timeUnit['unit']}&aggregate=${timeUnit['number']}");
  final res = await http.get(uri);

  if (res.statusCode == 200) {
    List<dynamic> data = jsonDecode(res.body);

    return data.map((e) {
      return Candle(
        date: DateTime.fromMillisecondsSinceEpoch(e[0]),
        open: roundToDecimals(e[1].toDouble(), decimals),
        high: roundToDecimals(e[2].toDouble(), decimals),
        low: roundToDecimals(e[3].toDouble(), decimals),
        close: roundToDecimals(e[4].toDouble(), decimals),
        volume: roundToDecimals(e[5].toDouble(), decimals),
      );
    }).toList();
  } else {
    throw Exception('Failed to load candles');
  }
}

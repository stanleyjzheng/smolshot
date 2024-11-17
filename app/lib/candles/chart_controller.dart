import 'package:fl_chart/fl_chart.dart';
import 'package:get/get.dart';
import 'package:k_chart/entity/k_line_entity.dart';

import './types.dart';
import './fetch_ohlc.dart';

class ChartController extends GetxController {
  static ChartController get to => Get.put(ChartController());

  List<KLineEntity> kChartCandles = <KLineEntity>[].obs;
  List<FlSpot> lineChart = <FlSpot>[].obs;

  String interval = '1';

  double inrRate = 0.0;

  getCandles({required Coin coinData, required String interval}) async {
    try {
      fetchCandles(
        symbol: coinData.id,
        interval: interval,
      ).then(
        (value) {
          kChartCandles.clear();
          lineChart.clear();
          for (int i = 0; i < value.length; i++) {
            kChartCandles.add(KLineEntity.fromCustom(
                time: value[i].date.millisecondsSinceEpoch,
                amount: value[i].high,
                change: value[i].high - value[i].low,
                close: value[i].close,
                high: value[i].high,
                low: value[i].low,
                open: value[i].open,
                vol: value[i].volume,
                ratio: value[i].low));

            lineChart.add(
              FlSpot(
                double.parse(value[i].date.millisecondsSinceEpoch.toString()),
                value[i].close,
              ),
            );
          }
        },
      );
    } catch (e) {
      kChartCandles.clear();
      lineChart.clear();
      return;
    }

    update();
  }

  updateCoinGraph(data, Coin coinData) {
    if (data.containsKey("k") == true &&
        kChartCandles[kChartCandles.length - 1].time! < data["k"]["t"]) {
      kChartCandles.add(KLineEntity.fromCustom(
          time: data["k"]["t"],
          amount: double.parse(data["k"]["h"].toString()),
          change: double.parse(data["k"]["v"].toString()),
          close: double.parse(data["k"]["c"].toString()),
          high: double.parse(data["k"]["h"].toString()),
          low: double.parse(data["k"]["l"].toString()),
          open: double.parse(data["k"]["o"].toString()),
          vol: double.parse(data["k"]["v"].toString()),
          ratio: double.parse(data["k"]["c"].toString())));
    } else if (data.containsKey("k") == true) {
      kChartCandles[kChartCandles.length - 1] = KLineEntity.fromCustom(
          time: data["k"]["t"],
          amount: double.parse(data["k"]["h"].toString()),
          change: double.parse(data["k"]["v"].toString()),
          close: double.parse(data["k"]["c"].toString()),
          high: double.parse(data["k"]["h"].toString()),
          low: double.parse(data["k"]["l"].toString()),
          open: double.parse(data["k"]["o"].toString()),
          vol: double.parse(data["k"]["v"].toString()),
          ratio: double.parse(data["k"]["c"].toString()));
    }

    update();
  }
}

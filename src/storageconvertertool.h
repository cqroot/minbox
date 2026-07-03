#ifndef STORAGECONVERTERTOOL_H
#define STORAGECONVERTERTOOL_H

#include <QWidget>

class QDoubleSpinBox;
class QComboBox;
class QLabel;

class StorageConverterTool : public QWidget {
    Q_OBJECT

  public:
    explicit StorageConverterTool(QWidget *parent = nullptr);
    ~StorageConverterTool();

  private slots:
    void onValueChanged(double value);
    void onUnitChanged(int index);

  private:
    void convertStorage(double value, int unitIndex);

    QDoubleSpinBox *m_value;
    QComboBox *m_unitCombo;
    QLabel *m_resultBytes;
    QLabel *m_resultKB;
    QLabel *m_resultMB;
    QLabel *m_resultGB;
    QLabel *m_resultTB;
    QLabel *m_resultKiB;
    QLabel *m_resultMiB;
    QLabel *m_resultGiB;
    QLabel *m_resultTiB;
};

#endif

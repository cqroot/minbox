#include "storageconvertertool.h"

#include <QDoubleSpinBox>
#include <QComboBox>
#include <QLabel>
#include <QGroupBox>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QFormLayout>
#include <QSpacerItem>

StorageConverterTool::StorageConverterTool(QWidget *parent)
    : QWidget(parent)
    , m_value(nullptr)
    , m_unitCombo(nullptr)
    , m_resultBytes(nullptr)
    , m_resultKB(nullptr)
    , m_resultMB(nullptr)
    , m_resultGB(nullptr)
    , m_resultTB(nullptr)
    , m_resultKiB(nullptr)
    , m_resultMiB(nullptr)
    , m_resultGiB(nullptr)
    , m_resultTiB(nullptr)
{
    setSizePolicy(QSizePolicy::Minimum, QSizePolicy::Minimum);

    QVBoxLayout *layout = new QVBoxLayout(this);

    layout->addStretch();

    QLabel *title = new QLabel("Storage Unit Conversion");
    title->setStyleSheet("font-size: 18px; font-weight: bold;");
    layout->addWidget(title);

    QHBoxLayout *inputLayout = new QHBoxLayout();
    QLabel *inputLabel = new QLabel("Value:");
    inputLayout->addWidget(inputLabel);

    m_value = new QDoubleSpinBox();
    m_value->setMinimum(0);
    m_value->setMaximum(999999999999);
    m_value->setValue(1);
    inputLayout->addWidget(m_value);

    m_unitCombo = new QComboBox();
    m_unitCombo->addItem("Bytes (B)");
    m_unitCombo->addItem("Kilobytes (KB)");
    m_unitCombo->addItem("Megabytes (MB)");
    m_unitCombo->addItem("Gigabytes (GB)");
    m_unitCombo->addItem("Terabytes (TB)");
    m_unitCombo->addItem("Petabytes (PB)");
    m_unitCombo->addItem("Kibibytes (KiB)");
    m_unitCombo->addItem("Mebibytes (MiB)");
    m_unitCombo->addItem("Gibibytes (GiB)");
    m_unitCombo->addItem("Tebibytes (TiB)");
    inputLayout->addWidget(m_unitCombo);

    layout->addLayout(inputLayout);

    QGroupBox *resultGroup = new QGroupBox("Conversion Result");
    QFormLayout *formLayout = new QFormLayout();

    m_resultBytes = new QLabel("1");
    formLayout->addRow("Bytes (B):", m_resultBytes);

    m_resultKB = new QLabel("0.0009765625");
    formLayout->addRow("Kilobytes (KB):", m_resultKB);

    m_resultMB = new QLabel("9.53674316e-7");
    formLayout->addRow("Megabytes (MB):", m_resultMB);

    m_resultGB = new QLabel("9.31322575e-10");
    formLayout->addRow("Gigabytes (GB):", m_resultGB);

    m_resultTB = new QLabel("9.09494702e-13");
    formLayout->addRow("Terabytes (TB):", m_resultTB);

    m_resultKiB = new QLabel("0.0009765625");
    formLayout->addRow("Kibibytes (KiB):", m_resultKiB);

    m_resultMiB = new QLabel("9.53674316e-7");
    formLayout->addRow("Mebibytes (MiB):", m_resultMiB);

    m_resultGiB = new QLabel("9.31322575e-10");
    formLayout->addRow("Gibibytes (GiB):", m_resultGiB);

    m_resultTiB = new QLabel("9.09494702e-13");
    formLayout->addRow("Tebibytes (TiB):", m_resultTiB);

    resultGroup->setLayout(formLayout);
    layout->addWidget(resultGroup);

    layout->addStretch();

    connect(m_value, QOverload<double>::of(&QDoubleSpinBox::valueChanged), this, &StorageConverterTool::onValueChanged);
    connect(m_unitCombo, QOverload<int>::of(&QComboBox::currentIndexChanged), this, &StorageConverterTool::onUnitChanged);

    convertStorage(m_value->value(), 0);
}

StorageConverterTool::~StorageConverterTool()
{
}

void StorageConverterTool::onValueChanged(double value)
{
    convertStorage(value, m_unitCombo->currentIndex());
}

void StorageConverterTool::onUnitChanged(int index)
{
    convertStorage(m_value->value(), index);
}

void StorageConverterTool::convertStorage(double value, int unitIndex)
{
    const double kib = 1024.0;
    const double mebib = kib * 1024.0;
    const double gibib = mebib * 1024.0;
    const double tebib = gibib * 1024.0;

    const double kb = 1000.0;
    const double mb = kb * 1000.0;
    const double gb = mb * 1000.0;
    const double tb = gb * 1000.0;
    const double pb = tb * 1000.0;

    double bytes = 0.0;
    switch (unitIndex) {
    case 0:
        bytes = value;
        break;
    case 1:
        bytes = value * kb;
        break;
    case 2:
        bytes = value * mb;
        break;
    case 3:
        bytes = value * gb;
        break;
    case 4:
        bytes = value * tb;
        break;
    case 5:
        bytes = value * pb;
        break;
    case 6:
        bytes = value * kib;
        break;
    case 7:
        bytes = value * mebib;
        break;
    case 8:
        bytes = value * gibib;
        break;
    case 9:
        bytes = value * tebib;
        break;
    }

    m_resultBytes->setText(QString::number(bytes, 'g', 15));
    m_resultKB->setText(QString::number(bytes / kb, 'g', 15));
    m_resultMB->setText(QString::number(bytes / mb, 'g', 15));
    m_resultGB->setText(QString::number(bytes / gb, 'g', 15));
    m_resultTB->setText(QString::number(bytes / tb, 'g', 15));

    m_resultKiB->setText(QString::number(bytes / kib, 'g', 15));
    m_resultMiB->setText(QString::number(bytes / mebib, 'g', 15));
    m_resultGiB->setText(QString::number(bytes / gibib, 'g', 15));
    m_resultTiB->setText(QString::number(bytes / tebib, 'g', 15));
}

#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QStringList>

class QLineEdit;
class QListWidget;
class QStackedWidget;

class BaseConverterTool;
class StorageConverterTool;

class MainWindow : public QMainWindow {
    Q_OBJECT

  public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

  private slots:
    void onToolListClicked(int row);
    void onSearchTextChanged(const QString &text);

  private:
    void setupUi();

    QWidget *m_centralWidget;
    QLineEdit *m_searchInput;
    QListWidget *m_toolList;
    QStackedWidget *m_contentStack;

    BaseConverterTool *m_baseConverter;
    StorageConverterTool *m_storageConverter;

    QStringList m_allTools;
};

#endif
